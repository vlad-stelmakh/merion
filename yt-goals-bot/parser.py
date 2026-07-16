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
    r"(?:^|\b)(?:1\s*[-.]?\s*(?:й|й)?\s*тайм|первый\s+тайм|1\s*т)\b",
    re.IGNORECASE,
)
SECOND_HALF_MARKER_RE = re.compile(
    r"(?:^|\b)(?:2\s*[-.]?\s*(?:й|й)?\s*тайм|второй\s+тайм|2\s*т)\b",
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


def _goals_from_minute_first_line(line: str, half: Half) -> list[Goal]:
    match = MINUTE_FIRST_RE.match(line)
    if not match:
        return []
    minute = int(match.group("minute"))
    scorer = match.group("scorer").strip()
    goal_half = Half.SECOND if minute > 45 else half
    return [Goal(minute=minute, scorer=scorer, half=goal_half)]


def _goals_from_scorer_line(line: str, half: Half) -> list[Goal]:
    minutes = [int(m) for m in MINUTE_TOKEN_RE.findall(line)]
    if not minutes:
        return []

    scorer = MINUTE_TOKEN_RE.sub("", line)
    scorer = re.sub(r"\s*,\s*", " ", scorer)
    scorer = re.sub(r"\s+", " ", scorer).strip(" ,")
    if not scorer:
        return []

    goals: list[Goal] = []
    for minute in minutes:
        goal_half = Half.SECOND if minute > 45 else half
        goals.append(Goal(minute=minute, scorer=scorer, half=goal_half))
    return goals


def parse_match_results(text: str) -> MatchResult:
    """Разбирает текст с результатами и списком голов."""
    text = strip_urls_and_commands(text)

    home_team: str | None = None
    away_team: str | None = None
    score_home: int | None = None
    score_away: int | None = None
    goals: list[Goal] = []

    blocks = [block.strip() for block in re.split(r"\n\s*\n", text) if block.strip()]
    goal_block_index = 0

    for block in blocks:
        lines = [line.strip() for line in block.splitlines() if line.strip()]
        if not lines:
            continue

        score = _parse_score_line(lines[0])
        if score and len(lines) == 1:
            home_team, away_team, score_home, score_away = score
            continue

        if score:
            home_team, away_team, score_home, score_away = score
            lines = lines[1:]

        side = "home" if goal_block_index == 0 else "away"
        goal_block_index += 1

        current_half = Half.FIRST
        seen_scorer_style = False
        current_side = side

        for line in lines:
            if SECOND_HALF_MARKER_RE.search(line):
                current_half = Half.SECOND
                continue
            if HALF_MARKER_RE.search(line) and not SECOND_HALF_MARKER_RE.search(line):
                current_half = Half.FIRST
                continue

            if _parse_score_line(line):
                continue

            minute_goals = _goals_from_minute_first_line(line, current_half)
            if minute_goals:
                for goal in minute_goals:
                    goal.side = current_side
                    goals.append(goal)
                continue

            scorer_goals = _goals_from_scorer_line(line, current_half)
            if scorer_goals:
                if not seen_scorer_style and goals:
                    current_side = "away"
                    current_half = Half.SECOND
                seen_scorer_style = True
                for goal in scorer_goals:
                    if goal.minute <= 45 and current_half == Half.SECOND and goal.half != Half.SECOND:
                        goal.half = Half.SECOND
                    goal.side = current_side
                    goals.append(goal)

    return MatchResult(
        home_team=home_team,
        away_team=away_team,
        score_home=score_home,
        score_away=score_away,
        goals=goals,
    )
