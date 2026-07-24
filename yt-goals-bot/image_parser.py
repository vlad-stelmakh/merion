"""Парсинг картинки с результатом турнира (OCR)."""

from __future__ import annotations

import io
import re
from functools import lru_cache

import easyocr
import numpy as np
from PIL import Image, ImageEnhance, ImageOps

from generator import HALF_MINUTES
from parser import Goal, Half, MatchResult

GOAL_CHUNK_RE = re.compile(
    r"(\d{1,2})\s*['′`´\"]+\s*(.+?)(?=\s*\d{1,2}\s*['′`´\"]|$)",
    re.DOTALL,
)
SCORE_RE = re.compile(r"\b(\d{1,2})\s*[-–—]\s*(\d{1,2})\b")
TEAM_RE = re.compile(r"^(?:FC|YFC|ФК)\s+[A-Za-zА-Яа-яЁё0-9][\w\s\-]*$", re.IGNORECASE)
BALL_CHARS = "⚽●◦·•°|"
NOISE_RE = re.compile(
    r"friends|f\.?f\.?f|for\s*friends|fnk|kucha$|^\d+\s*kucha",
    re.IGNORECASE,
)


@lru_cache(maxsize=1)
def _reader() -> easyocr.Reader:
    return easyocr.Reader(["ru", "en"], gpu=False, verbose=False)


def _preprocess_crop(image: Image.Image, box: tuple[int, int, int, int]) -> np.ndarray:
    crop = image.crop(box)
    crop = crop.resize((crop.width * 2, crop.height * 2), Image.Resampling.LANCZOS)
    crop = ImageOps.grayscale(crop)
    crop = ImageEnhance.Contrast(crop).enhance(1.8)
    return np.array(crop.convert("RGB"))


def _ocr_region(image: Image.Image, box: tuple[int, int, int, int]) -> str:
    arr = _preprocess_crop(image, box)
    parts = _reader().readtext(
        arr,
        paragraph=True,
        min_size=8,
        contrast_ths=0.05,
        adjust_contrast=0.7,
        text_threshold=0.6,
        low_text=0.3,
    )
    texts = [text for _box, text, conf in parts if conf > 0.15 and text.strip()]
    return "\n".join(texts)


def _clean_scorer(text: str) -> str:
    text = text.strip()
    for ch in BALL_CHARS:
        text = text.replace(ch, " ")
    text = re.sub(r"\s+", " ", text).strip(" -—.,\"")
    # убрать хвосты следующего гола
    text = re.sub(r"\s+\d{1,2}\s*['′`´\"].*$", "", text)
    return text.strip()


def _is_valid_scorer(name: str) -> bool:
    if not name or len(name) < 4:
        return False
    if NOISE_RE.search(name):
        return False
    if re.fullmatch(r"[\d\s\-]+", name):
        return False

    words = name.split()
    single_letter = sum(1 for w in words if len(w) == 1)
    if single_letter >= 2:
        return False
    if len(words) >= 2 and sum(len(w) <= 2 for w in words) >= len(words) - 1:
        return False
    if len(words) == 1 and len(words[0]) <= 3:
        return False

    letters = sum(ch.isalpha() for ch in name)
    return letters / len(name) >= 0.55


def _parse_goals_from_text(text: str, side: str) -> list[Goal]:
    goals: list[Goal] = []
    text = text.replace("\n", " ")
    for match in GOAL_CHUNK_RE.finditer(text):
        minute = int(match.group(1))
        if minute > 90:
            continue
        scorer = _clean_scorer(match.group(2))
        if not _is_valid_scorer(scorer):
            continue
        goals.append(
            Goal(
                minute=minute,
                scorer=scorer,
                half=_half_from_minute(minute),
                side=side,
            )
        )
    return goals


def _half_from_minute(minute: int) -> Half:
    return Half.SECOND if minute > HALF_MINUTES else Half.FIRST


def _parse_score(image: Image.Image) -> tuple[int, int] | None:
    w, h = image.size
    box = (int(w * 0.35), int(h * 0.22), int(w * 0.65), int(h * 0.42))
    text = _ocr_region(image, box)
    match = SCORE_RE.search(text.replace(" ", ""))
    if not match:
        match = SCORE_RE.search(text)
    if match:
        return int(match.group(1)), int(match.group(2))
    return None


def _parse_team_name(image: Image.Image, side: str) -> str | None:
    w, h = image.size
    if side == "home":
        box = (int(w * 0.02), int(h * 0.38), int(w * 0.48), int(h * 0.52))
    else:
        box = (int(w * 0.52), int(h * 0.38), int(w * 0.98), int(h * 0.52))

    text = _ocr_region(image, box)
    for line in text.splitlines():
        line = line.strip()
        if not line or NOISE_RE.search(line):
            continue
        if TEAM_RE.match(line):
            return line
        # FC Svyst / FC Kucha без префикса ФК
        if re.match(r"^[A-Za-zА-Яа-яЁё][\w\s\-]{2,}$", line) and len(line) <= 30:
            if side == "home" and "kucha" not in line.lower():
                return line if line.upper().startswith("FC") else f"FC {line}"
            if side == "away" and "kucha" in line.lower():
                return line if line.upper().startswith("FC") else f"FC {line}"
    return None


def _parse_goal_column(image: Image.Image, side: str) -> list[Goal]:
    w, h = image.size
    if side == "home":
        box = (int(w * 0.0), int(h * 0.50), int(w * 0.49), int(h * 0.98))
    else:
        box = (int(w * 0.51), int(h * 0.50), int(w * 0.99), int(h * 0.98))

    text = _ocr_region(image, box)
    return _parse_goals_from_text(text, side)


def _validate_match(match: MatchResult) -> None:
    if not match.goals:
        raise ValueError("Не удалось распознать голы на картинке.")

    if match.score_home is not None and match.score_away is not None:
        expected = match.score_home + match.score_away
        if len(match.goals) < expected:
            raise ValueError(
                f"Распознано только {len(match.goals)} из {expected} голов. "
                "Пришли фото чётче или добавь подпись с голами к картинке."
            )

    for goal in match.goals:
        if goal.minute <= 0:
            raise ValueError("Некорректные минуты голов — попробуй другое фото.")


def parse_tournament_image(image_bytes: bytes, caption: str | None = None) -> MatchResult:
    """Читает турнирную картинку и возвращает MatchResult."""
    if caption and caption.strip():
        from parser import parse_match_results

        captioned = parse_match_results(caption)
        if captioned.goals:
            return captioned

    image = Image.open(io.BytesIO(image_bytes)).convert("RGB")

    home_goals = _parse_goal_column(image, "home")
    away_goals = _parse_goal_column(image, "away")
    goals = sorted(home_goals + away_goals, key=lambda g: (g.half.value, g.minute))

    score = _parse_score(image)
    score_home, score_away = score if score else (None, None)

    home_team = _parse_team_name(image, "home")
    away_team = _parse_team_name(image, "away")

    if score_home is None:
        score_home = len(home_goals) or None
        score_away = len(away_goals) or None

    match = MatchResult(
        home_team=home_team,
        away_team=away_team,
        score_home=score_home,
        score_away=score_away,
        goals=goals,
    )
    _validate_match(match)
    return match
