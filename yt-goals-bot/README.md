# YouTube Goals Comment Bot

Telegram-бот, который принимает текст с результатами футбольного матча и две ссылки на записи 1-го и 2-го тайма на YouTube, а в ответ отдаёт готовый комментарий с таймкодами голов (минус 10 секунд от минуты гола).

## Быстрый старт

1. Создай бота у [@BotFather](https://t.me/BotFather) и скопируй токен.
2. Установи зависимости:

```bash
cd yt-goals-bot
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

3. Создай файл `.env`:

```bash
cp .env.example .env
# впиши BOT_TOKEN=...
```

4. Запусти бота:

```bash
python bot.py
```

## Как пользоваться

### Одним сообщением

```
FC Kucha 3:6 FC Serega United

15' Filipp PlusMinus (Sandro Karchava)
22' Filipp PlusMinus (Dima Semin)
40' Bobeeo — Sandro Karchava (автогол)

Evgensky 11', 20', 31'
Andrei Tereshkov 33'
Egor Levin 20', 49' (Levan Kviki)

https://www.youtube.com/watch?v=VIDEO_1ST_HALF
https://www.youtube.com/watch?v=VIDEO_2ND_HALF
```

### Пошагово

1. `/new`
2. Текст с результатами
3. Ссылка на 1-й тайм
4. Ссылка на 2-й тайм

## Логика таймкодов

- Каждый гол: **минута × 60 − 10 секунд**
- Голы в формате `15' Игрок` до 45-й минуты → 1-е видео
- Блок с голами вида `Игрок 11', 20'` → 2-е видео
- Минуты больше 45 (например `49'`) → 2-е видео

## Тесты

```bash
python -m unittest test_bot.py -v
```

## Команды бота

| Команда | Описание |
|---------|----------|
| `/start` | Приветствие |
| `/new` | Новый матч |
| `/help` | Инструкция и пример |
| `/cancel` | Отмена текущего ввода |
