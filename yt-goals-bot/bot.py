"""Telegram-бот: картинка турнира → ссылки → таймкоды YouTube."""

from __future__ import annotations

import asyncio
import logging
import os
import sys

from aiogram import Bot, Dispatcher, F
from aiogram.filters import Command, CommandStart
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.storage.memory import MemoryStorage
from aiogram.enums import ParseMode
from aiogram.types import BufferedInputFile, Message
from dotenv import load_dotenv

from image_parser import parse_tournament_image
from output import MatchOutput, build_comment_output
from parser import Goal, Half, MatchResult, extract_youtube_urls, parse_match_results

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

HELP_TEXT = """Привет! Я делаю <b>таймкоды YouTube</b> по картинке с результатом матча.

<b>Как пользоваться</b>
1. /new
2. Пришли <b>картинку</b> из группы турнира
   (лучше с <b>подписью</b> — список голов текстом)
3. Пришли <b>2 ссылки</b> на 1-й и 2-й тайм

<b>Ответ</b> — готовый текст с таймкодами для комментария YouTube.

Таймы по 25 минут: гол 41' → <code>15:50</code> во 2-м видео.
/cancel — отмена
"""


class MatchForm(StatesGroup):
    waiting_image = State()
    waiting_video_urls = State()


def _match_to_state(match: MatchResult) -> dict:
    return {
        "home_team": match.home_team,
        "away_team": match.away_team,
        "score_home": match.score_home,
        "score_away": match.score_away,
        "goals": [
            {
                "minute": g.minute,
                "scorer": g.scorer,
                "half": g.half.value,
                "side": g.side,
            }
            for g in match.goals
        ],
    }


def _match_from_state(data: dict) -> MatchResult:
    return MatchResult(
        home_team=data.get("home_team"),
        away_team=data.get("away_team"),
        score_home=data.get("score_home"),
        score_away=data.get("score_away"),
        goals=[
            Goal(
                minute=g["minute"],
                scorer=g["scorer"],
                half=Half(g["half"]),
                side=g.get("side"),
            )
            for g in data.get("goals", [])
        ],
    )


def _summary(match: MatchResult) -> str:
    teams = ""
    if match.home_team and match.away_team:
        sh = match.score_home if match.score_home is not None else "?"
        sa = match.score_away if match.score_away is not None else "?"
        teams = f"\n<b>{match.home_team}</b> {sh}:{sa} <b>{match.away_team}</b>"
    first = len([g for g in match.goals if g.half == Half.FIRST])
    second = len([g for g in match.goals if g.half == Half.SECOND])
    return f"Нашёл <b>{len(match.goals)}</b> голов ({first} в 1-м, {second} во 2-м тайме).{teams}"


async def _send_result(message: Message, result: MatchOutput) -> None:
    if result.image_png:
        photo = BufferedInputFile(result.image_png, filename="match_result.png")
        await message.answer_photo(photo=photo, caption=result.caption[:1024])
    else:
        await message.answer(result.caption)


async def cmd_start(message: Message, state: FSMContext) -> None:
    await state.clear()
    await message.answer("Привет! /new — начать. /help — инструкция.")


async def cmd_help(message: Message) -> None:
    await message.answer(HELP_TEXT, parse_mode=ParseMode.HTML)


async def cmd_new(message: Message, state: FSMContext) -> None:
    await state.clear()
    await state.set_state(MatchForm.waiting_image)
    await message.answer(
        "Пришли <b>картинку с результатом</b> из группы турнира.\n\n"
        "Можно добавить <b>подпись</b> с голами — тогда разбор будет точнее.",
        parse_mode=ParseMode.HTML,
    )


async def cmd_cancel(message: Message, state: FSMContext) -> None:
    await state.clear()
    await message.answer("Отменено. /new — начать заново.")


async def handle_image(message: Message, state: FSMContext, bot: Bot) -> None:
    if message.photo:
        file_id = message.photo[-1].file_id
    elif message.document and (message.document.mime_type or "").startswith("image/"):
        file_id = message.document.file_id
    else:
        await message.answer("Нужна картинка (фото). /help — инструкция.")
        return

    wait = await message.answer("Читаю картинку…")

    try:
        file = await bot.get_file(file_id)
        if not file.file_path:
            raise ValueError("Не удалось скачать файл.")
        downloaded = await bot.download_file(file.file_path)
        image_bytes = downloaded.read()
        caption = message.caption or ""
        match = await asyncio.to_thread(parse_tournament_image, image_bytes, caption or None)
    except Exception as exc:
        logger.exception("OCR failed")
        await wait.edit_text(f"Не смог разобрать картинку: {exc}\n\nПопробуй другое фото или /cancel.")
        return

    await state.update_data(match=_match_to_state(match))
    await state.set_state(MatchForm.waiting_video_urls)
    await wait.edit_text(
        f"{_summary(match)}\n\n"
        "Теперь пришли <b>2 ссылки</b> на YouTube (1-й и 2-й тайм) одним сообщением.",
        parse_mode=ParseMode.HTML,
    )


async def handle_video_urls(message: Message, state: FSMContext) -> None:
    text = message.text or ""
    urls = extract_youtube_urls(text)
    if len(urls) < 2:
        await message.answer(
            "Нужны 2 ссылки на YouTube.\nПример:\n"
            "https://youtu.be/VIDEO1\n"
            "https://youtu.be/VIDEO2",
        )
        return

    data = await state.get_data()
    match = _match_from_state(data.get("match", {}))

    try:
        result = build_comment_output(match, urls[0], urls[1])
    except ValueError as exc:
        await message.answer(f"Ошибка: {exc}")
        return

    await state.clear()
    await _send_result(message, result)


async def handle_outside_flow(message: Message, state: FSMContext) -> None:
    await message.answer("Отправь /new чтобы начать, или /help для инструкции.")


async def main() -> None:
    token = os.getenv("BOT_TOKEN")
    if not token:
        logger.error("Укажи BOT_TOKEN в .env")
        sys.exit(1)

    bot = Bot(token=token)
    dp = Dispatcher(storage=MemoryStorage())

    dp.message.register(cmd_start, CommandStart())
    dp.message.register(cmd_help, Command("help"))
    dp.message.register(cmd_new, Command("new"))
    dp.message.register(cmd_cancel, Command("cancel"))
    dp.message.register(
        handle_image,
        MatchForm.waiting_image,
        F.photo | F.document,
    )
    dp.message.register(
        lambda m, s: m.answer("Сначала пришли картинку с результатом матча."),
        MatchForm.waiting_image,
        F.text,
    )
    dp.message.register(
        handle_video_urls,
        MatchForm.waiting_video_urls,
        F.text,
    )
    dp.message.register(handle_outside_flow, F.text)

    logger.info("Бот запущен")
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
