"""Парсер текста с результатами футбольного матча."""

from __future__ import annotations

import re
from dataclasses import dataclass
from enum import Enum


class Half(Enum):
    FIRST = 1
    SECOND = 2


@dataclass
class Goal:
    minute: int
    scorer: str
    half: Half
    side: str | None = None  # "home" | "away"


@dataclass
class MatchResult:
    home_team: str | None
    away_team: str | None
    score_home: int | None
    score_away: int | None
    goals: list[Goal]


SCORE_RE = re.compile(
    r"^\s*(?P<home>.+?)\s+(?P<sh>\d+)\s*[-:–—]\s*(?P<sa>\d+)\s+(?P<away>.+?)\s*$",
    re.IGNORECASE,
)
MINUTE_FIRST_RE = re.compile(
    r"^\s*(?P<minute>\d{1,3})\s*['′`]\s*[-—–]?\s*(?P<scorer>.+?)\s*$"
)
MINUTE_TOKEN_RE = re.compile(r"(\d{1,3})\s*['′`]")
HALF_MARKER_RE = re.compile(
    r"^\s*(?:1\s*[-.]?\s*й\s*тайм|первый\s+тайм|1\s*т)\s*$",
    re.IGNORECASE,
)
SECOND_HALF_MARKER_RE = re.compile(
    r"^\s*(?:2\s*[-.]?\s*й\s*тайм|второй\s+тайм|2\s*т)\s*$",
    re.IGNORECASE,
)
YOUTUBE_URL_RE = re.compile(
    r"https?://(?:www\.)?(?:youtube\.com/watch\?[^\s]+|youtu\.be/[^\s]+)",
    re.IGNORECASE,
)


def extract_youtube_urls(text: str) -> list[str]:
    return YOUTUBE_URL_RE.findall(text)


def strip_urls_and_commands(text: str) -> str:
    text = YOUTUBE_URL_RE.sub("", text)
    text = re.sub(r"^\s*#\s*[12]\s*$", "", text, flags=re.MULTILINE)
    return text.strip()


def _parse_score_line(line: str) -> tuple[str, str, int, int] | None:
    match = SCORE_RE.match(line)
    if not match:
        return None
    return (
        match.group("home").strip(),
        match.group("away").strip(),
        int(match.group("sh")),
        int(match.group("sa")),
    )


def _resolve_half(minute: int, current_half: Half) -> Half:
    if minute > 45:
        return Half.SECOND
    return current_half


def _goals_from_minute_first_line(line: str, half: Half) -> list[Goal]:
    match = MINUTE_FIRST_RE.match(line)
    if not match:
        return []
    minute = int(match.group("minute"))
    scorer = match.group("scorer").strip()
    return [Goal(minute=minute, scorer=scorer, half=_resolve_half(minute, half))]


def _goals_from_scorer_line(line: str, half: Half) -> list[Goal]:
    minutes = [int(m) for m in MINUTE_TOKEN_RE.findall(line)]
    if not minutes:
        return []

    scorer = MINUTE_TOKEN_RE.sub("", line)
    scorer = re.sub(r"\s*,\s*", " ", scorer)
    scorer = re.sub(r"\s+", " ", scorer).strip(" ,")
    if not scorer:
        return []

    return [
        Goal(minute=minute, scorer=scorer, half=_resolve_half(minute, half))
        for minute in minutes
    ]


def _is_half_marker(line: str) -> Half | None:
    if SECOND_HALF_MARKER_RE.match(line):
        return Half.SECOND
    if HALF_MARKER_RE.match(line):
        return Half.FIRST
    return None


def _block_has_goals(block: str) -> bool:
    for line in block.splitlines():
        line = line.strip()
        if not line or _is_half_marker(line) or _parse_score_line(line):
            continue
        if _goals_from_minute_first_line(line, Half.FIRST) or _goals_from_scorer_line(line, Half.FIRST):
            return True
    return False


def parse_match_results(text: str) -> MatchResult:
    """Разбирает текст с результатами и списком голов."""
    text = strip_urls_and_commands(text)

    home_team: str | None = None
    away_team: str | None = None
    score_home: int | None = None
    score_away: int | None = None
    goals: list[Goal] = []

    has_half_markers = bool(
        HALF_MARKER_RE.search(text, re.MULTILINE) or SECOND_HALF_MARKER_RE.search(text, re.MULTILINE)
    )

    # Разделение команд по пустым строкам (если нет явных таймов)
    goal_blocks = [b.strip() for b in re.split(r"\n\s*\n", text) if b.strip() and _block_has_goals(b)]
    goal_block_sides: dict[int, str] = {}
    if not has_half_markers and goal_blocks:
        for idx in range(len(goal_blocks)):
            goal_block_sides[idx] = "home" if idx == 0 else "away"

    current_half = Half.FIRST
    seen_scorer_style = False
    current_side = "home"
    goal_block_index = 0
    in_goal_block = False

    for raw_line in text.splitlines():
        line = raw_line.strip()
        if not line:
            if in_goal_block:
                goal_block_index += 1
                in_goal_block = False
            continue

        half_marker = _is_half_marker(line)
        if half_marker is not None:
            current_half = half_marker
            continue

        score = _parse_score_line(line)
        if score:
            home_team, away_team, score_home, score_away = score
            continue

        if not in_goal_block:
            in_goal_block = True
            if not has_half_markers and goal_block_index in goal_block_sides:
                current_side = goal_block_sides[goal_block_index]

        minute_goals = _goals_from_minute_first_line(line, current_half)
        if minute_goals:
            for goal in minute_goals:
                goal.side = None if has_half_markers else current_side
                goals.append(goal)
            continue

        scorer_goals = _goals_from_scorer_line(line, current_half)
        if scorer_goals:
            if not seen_scorer_style and goals and not has_half_markers:
                current_side = "away"
                current_half = Half.SECOND
            seen_scorer_style = True
            for goal in scorer_goals:
                if goal.minute <= 45 and current_half == Half.SECOND:
                    goal.half = Half.SECOND
                goal.side = None if has_half_markers else current_side
                goals.append(goal)

    return MatchResult(
        home_team=home_team,
        away_team=away_team,
        score_home=score_home,
        score_away=score_away,
        goals=goals,
    )
