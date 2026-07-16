"""Telegram-бот для генерации комментария с таймкодами голов."""

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

from output import MatchOutput, build_match_output
from parser import extract_youtube_urls, parse_match_results

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

HELP_TEXT = """Привет! Я делаю <b>картинку с результатом</b> и подпись с таймкодами для YouTube.

<b>Быстрый способ</b> — одним сообщением:
1. Текст с результатами
2. Две ссылки на 1-й и 2-й тайм

<b>Пошаговый способ</b>
/new → текст → ссылка 1-го тайма → ссылка 2-го тайма

<b>Формат результатов</b>
<code>YFC Pakuta 6:2 City United FC

4' David Benidze (Vova Orange)
8' Vova Orange (David Benidze)
18' Luka Kaladze (Arthur Parkour)

6' Efim Tarasenko (Danya Shangin)
20' Efim Tarasenko (Danya Shangin)</code>

Между голами разных команд — пустая строка.
Таймкод в подписи = минута гола − 10 секунд.
"""


class MatchForm(StatesGroup):
    waiting_results = State()
    waiting_first_half_url = State()
    waiting_second_half_url = State()


def _assign_half_urls(urls: list[str]) -> tuple[str, str] | None:
    if len(urls) < 2:
        return None
    return urls[0], urls[1]


def _try_generate_from_single_message(text: str) -> MatchOutput | None:
    urls = extract_youtube_urls(text)
    if len(urls) < 2:
        return None

    match = parse_match_results(text)
    if not match.goals:
        return None

    first_url, second_url = urls[0], urls[1]
    return build_match_output(match, first_url, second_url)


async def _send_result(message: Message, result: MatchOutput) -> None:
    photo = BufferedInputFile(result.image_png, filename="match_result.png")
    await message.answer_photo(photo=photo, caption=result.caption)


async def cmd_start(message: Message, state: FSMContext) -> None:
    await state.clear()
    await message.answer(
        "Привет! Отправь /new чтобы начать, или /help для инструкции.",
    )


async def cmd_help(message: Message) -> None:
    await message.answer(HELP_TEXT, parse_mode=ParseMode.HTML)


async def cmd_new(message: Message, state: FSMContext) -> None:
    await state.clear()
    await state.set_state(MatchForm.waiting_results)
    await message.answer(
        "Пришли текст с результатами матча и голами.\n\n"
        "Можно сразу добавить 2 ссылки на таймы — тогда отвечу сразу.",
    )


async def cmd_cancel(message: Message, state: FSMContext) -> None:
    await state.clear()
    await message.answer("Отменено. /new — начать заново.")


async def handle_results(message: Message, state: FSMContext) -> None:
    text = message.text or ""
    ready = _try_generate_from_single_message(text)
    if ready:
        await state.clear()
        await _send_result(message, ready)
        return

    match = parse_match_results(text)
    if not match.goals:
        await message.answer(
            "Не нашёл голов в тексте. Проверь формат или отправь /help.",
        )
        return

    await state.update_data(results_text=text)
    await state.set_state(MatchForm.waiting_first_half_url)
    teams_hint = ""
    if match.home_team and match.away_team:
        teams_hint = f"\n\nМатч: {match.home_team} — {match.away_team}"
    await message.answer(
        f"Нашёл {len(match.goals)} голов.{teams_hint}\n\n"
        "Теперь пришли ссылку на запись <b>1-го тайма</b>.",
        parse_mode=ParseMode.HTML,
    )


async def handle_first_half_url(message: Message, state: FSMContext) -> None:
    text = message.text or ""
    urls = extract_youtube_urls(text)
    if not urls:
        await message.answer("Нужна ссылка на YouTube. Пример: https://youtu.be/XXXXXXXXXXX")
        return

    await state.update_data(first_half_url=urls[0])
    await state.set_state(MatchForm.waiting_second_half_url)
    await message.answer(
        "Отлично. Теперь ссылка на запись <b>2-го тайма</b>.",
        parse_mode=ParseMode.HTML,
    )


async def handle_second_half_url(message: Message, state: FSMContext) -> None:
    text = message.text or ""
    urls = extract_youtube_urls(text)
    if not urls:
        await message.answer("Нужна ссылка на YouTube.")
        return

    data = await state.get_data()
    results_text = data.get("results_text", "")
    first_half_url = data.get("first_half_url", "")

    try:
        match = parse_match_results(results_text)
        result = build_match_output(match, first_half_url, urls[0])
    except ValueError as exc:
        await message.answer(f"Ошибка: {exc}")
        return

    await state.clear()
    await _send_result(message, result)


async def handle_outside_flow(message: Message, state: FSMContext) -> None:
    text = message.text or ""
    ready = _try_generate_from_single_message(text)
    if ready:
        await _send_result(message, ready)
        return

    urls = extract_youtube_urls(text)
    if len(urls) >= 2:
        await message.answer(
            "Ссылки есть, но в тексте не найдены голы. Добавь результат матча или /new.",
        )
        return

    await message.answer("Отправь /new чтобы начать, или /help для примера.")


async def main() -> None:
    token = os.getenv("BOT_TOKEN")
    if not token:
        logger.error("Укажи BOT_TOKEN в .env или переменных окружения.")
        sys.exit(1)

    bot = Bot(token=token)
    dp = Dispatcher(storage=MemoryStorage())

    dp.message.register(cmd_start, CommandStart())
    dp.message.register(cmd_help, Command("help"))
    dp.message.register(cmd_new, Command("new"))
    dp.message.register(cmd_cancel, Command("cancel"))
    dp.message.register(
        handle_results,
        MatchForm.waiting_results,
        F.text,
    )
    dp.message.register(
        handle_first_half_url,
        MatchForm.waiting_first_half_url,
        F.text,
    )
    dp.message.register(
        handle_second_half_url,
        MatchForm.waiting_second_half_url,
        F.text,
    )
    dp.message.register(handle_outside_flow, F.text)

    logger.info("Бот запущен")
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
