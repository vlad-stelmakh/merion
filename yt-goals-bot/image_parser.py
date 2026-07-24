"""Парсинг картинки с результатом турнира (OCR)."""

from __future__ import annotations

import io
import re
from functools import lru_cache

import easyocr
import numpy as np
from PIL import Image

from generator import HALF_MINUTES
from parser import Goal, Half, MatchResult

GOAL_LINE_RE = re.compile(
    r"^(\d{1,3})\s*['′`´]?\s*(.+)$",
)
SCORE_RE = re.compile(r"(\d{1,2})\s*[-–—:]\s*(\d{1,2})")
TEAM_RE = re.compile(r"^(FC|YFC|ФК)\s+.+$", re.IGNORECASE)
BALL_CHARS = "⚽●◦·•oO0°"


@lru_cache(maxsize=1)
def _reader() -> easyocr.Reader:
    return easyocr.Reader(["ru", "en"], gpu=False, verbose=False)


def _clean_scorer(text: str) -> str:
    text = text.strip()
    for ch in BALL_CHARS:
        text = text.replace(ch, " ")
    text = re.sub(r"\s+", " ", text).strip(" -—")
    return text


def _parse_goal_text(text: str) -> tuple[int, str] | None:
    text = text.strip()
    match = GOAL_LINE_RE.match(text)
    if not match:
        return None
    scorer = _clean_scorer(match.group(2))
    if not scorer or scorer.isdigit():
        return None
    return int(match.group(1)), scorer


def _half_from_minute(minute: int) -> Half:
    return Half.SECOND if minute > HALF_MINUTES else Half.FIRST


def _merge_ocr_lines(items: list[tuple], y_threshold: int = 18) -> list[tuple[str, float, float, float, float]]:
    """Склеивает соседние OCR-фрагменты в строки."""
    if not items:
        return []

    sorted_items = sorted(items, key=lambda x: (x[0][0][1], x[0][0][0]))
    lines: list[tuple[str, float, float, float, float]] = []
    current_text = ""
    current_box = None

    for box, text, _conf in sorted_items:
        y = box[0][1]
        if current_box is None:
            current_text = text
            current_box = box
            continue

        if abs(y - current_box[0][1]) <= y_threshold:
            current_text += " " + text
            xs = [p[0] for p in current_box + box]
            ys = [p[1] for p in current_box + box]
            current_box = [
                [min(xs), min(ys)],
                [max(xs), min(ys)],
                [max(xs), max(ys)],
                [min(xs), max(ys)],
            ]
        else:
            cx = sum(p[0] for p in current_box) / 4
            cy = sum(p[1] for p in current_box) / 4
            lines.append((current_text.strip(), cx, cy, current_box[2][0] - current_box[0][0], current_box[2][1] - current_box[0][1]))
            current_text = text
            current_box = box

    if current_box is not None:
        cx = sum(p[0] for p in current_box) / 4
        cy = sum(p[1] for p in current_box) / 4
        lines.append((current_text.strip(), cx, cy, current_box[2][0] - current_box[0][0], current_box[2][1] - current_box[0][1]))

    return lines


def parse_tournament_image(image_bytes: bytes) -> MatchResult:
    """Читает турнирную картинку и возвращает MatchResult."""
    image = Image.open(io.BytesIO(image_bytes)).convert("RGB")
    width, height = image.size
    arr = np.array(image)

    raw = _reader().readtext(arr)
    lines = _merge_ocr_lines(raw)

    home_team: str | None = None
    away_team: str | None = None
    score_home: int | None = None
    score_away: int | None = None
    goals: list[Goal] = []

    # Счёт — в верхней половине, ближе к центру
    for text, cx, cy, _w, _h in lines:
        if cy > height * 0.55:
            continue
        score_match = SCORE_RE.search(text.replace(" ", ""))
        if score_match and abs(cx - width / 2) < width * 0.2:
            score_home = int(score_match.group(1))
            score_away = int(score_match.group(2))
            break

    # Названия команд
    team_candidates: list[tuple[str, float, float]] = []
    for text, cx, cy, _w, _h in lines:
        if cy > height * 0.62 or cy < height * 0.18:
            continue
        cleaned = text.strip()
        if TEAM_RE.match(cleaned) or (len(cleaned) > 3 and cleaned[0].isupper() and "'" not in cleaned):
            if not SCORE_RE.fullmatch(cleaned.replace(" ", "")):
                team_candidates.append((cleaned, cx, cy))

    team_candidates.sort(key=lambda item: item[1])
    if len(team_candidates) >= 2:
        home_team = team_candidates[0][0]
        away_team = team_candidates[-1][0]
    elif len(team_candidates) == 1:
        home_team = team_candidates[0][0]

    # Голы — нижняя часть, две колонки
    goal_zone_y = height * 0.45
    for text, cx, cy, _w, _h in lines:
        if cy < goal_zone_y:
            continue
        parsed = _parse_goal_text(text)
        if not parsed:
            continue
        minute, scorer = parsed
        side = "home" if cx < width / 2 else "away"
        goals.append(
            Goal(
                minute=minute,
                scorer=scorer,
                half=_half_from_minute(minute),
                side=side,
            )
        )

    goals.sort(key=lambda g: (g.half.value, g.minute))

    if not goals:
        raise ValueError(
            "Не удалось распознать голы на картинке. "
            "Пришли фото как изображение (не файл) и без сильного сжатия."
        )

    if score_home is None:
        home_count = len([g for g in goals if g.side == "home"])
        away_count = len([g for g in goals if g.side == "away"])
        score_home = home_count or None
        score_away = away_count or None

    return MatchResult(
        home_team=home_team,
        away_team=away_team,
        score_home=score_home,
        score_away=score_away,
        goals=goals,
    )
