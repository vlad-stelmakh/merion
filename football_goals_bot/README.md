# Football Goals Telegram Bot

Telegram-бот, который из текста результатов матча и двух ссылок на таймы собирает комментарий с таймкодами голов (−10 секунд до удара).

## Что делает

1. Принимает текст с голами по таймам
2. Принимает YouTube-ссылку на **1 тайм**
3. Принимает YouTube-ссылку на **2 тайм**
4. Отдаёт готовый комментарий: для каждого гола ссылка на нужное видео с `t=секунда−10`

## Быстрый старт

```bash
cd football_goals_bot
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt

cp .env.example .env
# впиши BOT_TOKEN от @BotFather

python bot.py
```

## Диалог с ботом

| Шаг | Что отправить |
|-----|----------------|
| `/new` | начать матч |
| текст результатов | команды + голы |
| ссылка YouTube | 1 тайм |
| ссылка YouTube | 2 тайм |

Команды: `/start` `/help` `/new` `/cancel`

## Формат результатов

```text
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
```

Поддерживаются маркеры: `1 тайм`, `первый тайм`, `2 тайм`, `второй тайм` (и english half).  
Время: `15'`, `15`, `15:30`.

## Пример вывода

```text
⚽ Голы матча FC Kucha — FC Serega United (3:6)

1 тайм:
14:50 https://youtu.be/VIDEO1?t=890 — Filipp PlusMinus (Sandro Karchava)
21:50 https://youtu.be/VIDEO1?t=1310 — Filipp PlusMinus (Dima Semin)
...

2 тайм:
10:50 https://youtu.be/VIDEO2?t=650 — Evgensky (Egor Levin)
...
```

Время на видео считается **внутри тайма** (как в записи тайма), не как сквозные минуты всего матча.

## Тесты

```bash
cd football_goals_bot
pytest -q
```

## Структура

```text
football_goals_bot/
  bot.py            # Telegram FSM
  match_parser.py   # разбор текста результатов
  comment.py        # генерация комментария
  youtube.py        # video id и таймкоды
  requirements.txt
  .env.example
  tests/
```
