# YouTube Goals Comment Bot

Telegram-бот, который принимает текст с результатами футбольного матча и две ссылки на записи 1-го и 2-го тайма на YouTube. В ответ отправляет **PNG-карточку** с результатом (как в лиге F.F.F.) и подпись с таймкодами для комментария (минус 10 секунд от минуты гола).

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
YFC Pakuta 6:2 City United FC

4' David Benidze (Vova Orange)
8' Vova Orange (David Benidze)
18' Luka Kaladze (Arthur Parkour)
23' Arthur Parkour (Luka Kaladze)
38' Arman (David Benidze)
45' Luka Kaladze (Arthur Parkour)

6' Efim Tarasenko (Danya Shangin)
20' Efim Tarasenko (Danya Shangin)

https://www.youtube.com/watch?v=VIDEO_1ST_HALF
https://www.youtube.com/watch?v=VIDEO_2ND_HALF
```

**Важно:** между голами разных команд — пустая строка.

### Пошагово

1. `/new`
2. Текст с результатами
3. Ссылка на 1-й тайм
4. Ссылка на 2-й тайм

## Что получаешь

1. **Картинка** — зелёная карточка со счётом, двумя колонками голов и минутами (`4' ⚽ Игрок (ассист)`)
2. **Подпись к фото** — текст с YouTube-ссылками и таймкодами (−10 сек)

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
