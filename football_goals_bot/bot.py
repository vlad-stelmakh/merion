"""
Telegram-бот: из текста результатов матча + 2 ссылок на таймы
собирает комментарий с таймкодами голов (−10 сек).

Запуск:
    export BOT_TOKEN=...
    python bot.py
"""

from __future__ import annotations

import asyncio
import logging
import os

from aiogram import Bot, Dispatcher, F
from aiogram.filters import Command, CommandStart
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from aiogram.types import Message
from dotenv import load_dotenv

from comment import generate_comment, generate_telegram_html
from match_parser import parse_results
from youtube import extract_video_id, is_youtube_url

load_dotenv()

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)
logger = logging.getLogger("football_goals_bot")

OFFSET_SEC = 10

HELP_TEXT = """\
⚽ <b>Бот таймкодов голов</b>

Я собираю комментарий для YouTube: таймкод гола минус 10 секунд, \
с ссылками на нужный тайм.

<b>Как пользоваться</b>
1. /new — начать новый матч
2. Пришли текст с результатами (голы по таймам)
3. Пришли ссылку на <b>1 тайм</b>
4. Пришли ссылку на <b>2 тайм</b>
5. Скопируй готовый комментарий

<b>Пример текста результатов</b>
<pre>
FC Kucha 3:6 FC Serega United

1 тайм:
15' Filipp PlusMinus (Sandro Karchava)
22' Filipp PlusMinus (Dima Semin)
40' Bobeeo — Sandro Karchava (автогол)

2 тайм:
11' Evgensky (Egor Levin)
20' Egor Levin
31' Evgensky
33' Andrei Tereshkov (Vlados Vlad)
40' Evgensky (Vlados Vlad)
49' Egor Levin (Levan Kviki)
</pre>

Команды: /new /cancel /help
"""


class Form(StatesGroup):
    waiting_results = State()
    waiting_half1 = State()
    waiting_half2 = State()


async def cmd_start(message: Message, state: FSMContext) -> None:
    await state.clear()
    await message.answer(HELP_TEXT, parse_mode="HTML")
    await message.answer(
        "Жду текст с результатами матча.\n"
        "Или нажми /new чтобы начать заново."
    )
    await state.set_state(Form.waiting_results)


async def cmd_help(message: Message) -> None:
    await message.answer(HELP_TEXT, parse_mode="HTML")


async def cmd_new(message: Message, state: FSMContext) -> None:
    await state.clear()
    await state.set_state(Form.waiting_results)
    await message.answer(
        "Ок, новый матч.\n\n"
        "Пришли текст с результатами — команды и голы по таймам "
        "(как в /help)."
    )


async def cmd_cancel(message: Message, state: FSMContext) -> None:
    await state.clear()
    await message.answer("Отменено. /new — начать снова.")


async def on_results(message: Message, state: FSMContext) -> None:
    text = (message.text or "").strip()
    if not text:
        await message.answer("Пришли текст с результатами.")
        return

    match = parse_results(text)
    if not match.goals:
        await message.answer(
            "Не нашёл голы в тексте.\n\n"
            "Нужны строки вида:\n"
            "<code>15' Имя игрока</code>\n"
            "и заголовки <code>1 тайм:</code> / <code>2 тайм:</code>",
            parse_mode="HTML",
        )
        return

    half1_count = sum(1 for g in match.goals if g.half == 1)
    half2_count = sum(1 for g in match.goals if g.half == 2)

    await state.update_data(results_text=text)
    await state.set_state(Form.waiting_half1)

    teams = f"{match.home_team} — {match.away_team}"
    score = f" ({match.score})" if match.score else ""
    await message.answer(
        f"Нашёл матч: <b>{teams}</b>{score}\n"
        f"Голов: {len(match.goals)} "
        f"(1 тайм: {half1_count}, 2 тайм: {half2_count})\n\n"
        "Теперь пришли <b>ссылку на 1 тайм</b> (YouTube).",
        parse_mode="HTML",
    )


async def on_half1(message: Message, state: FSMContext) -> None:
    url = (message.text or "").strip()
    video_id = extract_video_id(url)
    if not video_id:
        await message.answer(
            "Это не похоже на YouTube-ссылку.\n"
            "Пример: https://youtu.be/xxxxxxxxxxx\n"
            "или https://www.youtube.com/watch?v=xxxxxxxxxxx"
        )
        return

    await state.update_data(half1_video_id=video_id, half1_url=url)
    await state.set_state(Form.waiting_half2)
    await message.answer(
        "1 тайм принят ✅\nПришли <b>ссылку на 2 тайм</b>.",
        parse_mode="HTML",
    )


async def on_half2(message: Message, state: FSMContext) -> None:
    url = (message.text or "").strip()
    video_id = extract_video_id(url)
    if not video_id:
        await message.answer(
            "Это не похоже на YouTube-ссылку.\n"
            "Пример: https://youtu.be/xxxxxxxxxxx"
        )
        return

    data = await state.get_data()
    half1_id = data.get("half1_video_id")
    results_text = data.get("results_text", "")
    if not half1_id or not results_text:
        await state.clear()
        await message.answer("Сессия сброшена. Начни с /new")
        return

    match = parse_results(results_text)
    html = generate_telegram_html(
        match, half1_id, video_id, offset_sec=OFFSET_SEC
    )
    plain = generate_comment(
        match, half1_id, video_id, offset_sec=OFFSET_SEC
    )

    await message.answer(
        html + f"\n\n<i>Таймкоды: −{OFFSET_SEC} сек до гола</i>",
        parse_mode="HTML",
        disable_web_page_preview=True,
    )
    await message.answer(
        "📋 <b>Текст для копирования в YouTube:</b>\n\n"
        f"<pre>{_esc_pre(plain)}</pre>",
        parse_mode="HTML",
        disable_web_page_preview=True,
    )
    await state.clear()
    await message.answer("Готово. /new — следующий матч.")


async def on_unexpected(message: Message, state: FSMContext) -> None:
    current = await state.get_state()
    if current is None:
        if message.text and is_youtube_url(message.text):
            await message.answer("Сначала пришли результаты. Жми /new")
            return
        await message.answer("Начни с /new или посмотри /help")
        return
    await message.answer("Жду данные по шагам. /cancel — отмена, /help — справка.")


def _esc_pre(text: str) -> str:
    return (
        text.replace("&", "&amp;")
        .replace("<", "&lt;")
        .replace(">", "&gt;")
    )


def build_dispatcher() -> Dispatcher:
    dp = Dispatcher()
    dp.message.register(cmd_start, CommandStart())
    dp.message.register(cmd_help, Command("help"))
    dp.message.register(cmd_new, Command("new"))
    dp.message.register(cmd_cancel, Command("cancel"))
    dp.message.register(on_results, Form.waiting_results, F.text)
    dp.message.register(on_half1, Form.waiting_half1, F.text)
    dp.message.register(on_half2, Form.waiting_half2, F.text)
    dp.message.register(on_unexpected)
    return dp


async def main() -> None:
    token = os.getenv("BOT_TOKEN", "").strip()
    if not token:
        raise SystemExit(
            "Не задан BOT_TOKEN. Скопируй .env.example → .env "
            "или export BOT_TOKEN=..."
        )

    bot = Bot(token=token)
    dp = build_dispatcher()
    logger.info("Бот запускается (polling)...")
    await dp.start_polling(bot)


if __name__ == "__main__":
    asyncio.run(main())
