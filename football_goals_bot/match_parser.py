"""Парсер текста с результатами матча и голами по таймам."""

from __future__ import annotations

import re
from dataclasses import dataclass, field


HALF1_MARKERS = re.compile(
    r"(?i)^\s*(?:#{1,3}\s*)?(?:1(?:-й|ый)?\s*тайм|первый\s*тайм|1st\s*half|first\s*half)\s*:?\s*$"
)
HALF2_MARKERS = re.compile(
    r"(?i)^\s*(?:#{1,3}\s*)?(?:2(?:-й|ой)?\s*тайм|второй\s*тайм|2nd\s*half|second\s*half)\s*:?\s*$"
)

# 15' Name | 15′ Name | 15" Name | 15: Name | 15 - Name | 15. Name
GOAL_LINE = re.compile(
    r"^\s*"
    r"(?P<time>\d{1,3}(?::\d{2})?)\s*['′\"`´]?\s*"
    r"[-–—:]?\s*"
    r"(?P<scorer>.+?)\s*$"
)

SCORE_LINE = re.compile(
    r"(?P<home>.+?)\s+(\d+)\s*[:：-]\s*(\d+)\s+(?P<away>.+)"
)

VS_LINE = re.compile(
    r"(?P<home>.+?)\s+(?:—|-|vs\.?|против)\s+(?P<away>.+)",
    re.IGNORECASE,
)


@dataclass
class Goal:
    time_raw: str
    scorer: str
    half: int  # 1 или 2
    minutes: int
    seconds: int = 0

    @property
    def total_seconds(self) -> int:
        return self.minutes * 60 + self.seconds


@dataclass
class MatchResult:
    home_team: str = "Команда 1"
    away_team: str = "Команда 2"
    score: str | None = None
    goals: list[Goal] = field(default_factory=list)


def _parse_time(time_str: str) -> tuple[int, int]:
    """'15' / '15:30' -> (minutes, seconds)."""
    time_str = time_str.strip()
    if ":" in time_str:
        parts = time_str.split(":")
        return int(parts[0]), int(parts[1])
    return int(time_str), 0


def parse_results(text: str) -> MatchResult:
    """
    Разбирает текст результатов.

    Ожидаемый формат (гибкий):
        FC Kucha 3:6 FC Serega United

        1 тайм:
        15' Filipp PlusMinus (Sandro Karchava)
        22' Filipp PlusMinus (Dima Semin)

        2 тайм:
        11' Evgensky (Egor Levin)
        20' Egor Levin
    """
    result = MatchResult()
    current_half = 1
    lines = text.strip().splitlines()

    for raw_line in lines:
        line = raw_line.strip()
        if not line:
            continue

        if HALF1_MARKERS.match(line):
            current_half = 1
            continue
        if HALF2_MARKERS.match(line):
            current_half = 2
            continue

        score_match = SCORE_LINE.search(line)
        if score_match and not result.score:
            result.home_team = score_match.group("home").strip(" -—")
            result.away_team = score_match.group("away").strip(" -—")
            result.score = (
                f"{score_match.group(2)}:{score_match.group(3)}"
            )
            continue

        vs_match = VS_LINE.fullmatch(line)
        if vs_match and result.home_team == "Команда 1":
            result.home_team = vs_match.group("home").strip()
            result.away_team = vs_match.group("away").strip()
            continue

        goal_match = GOAL_LINE.match(line)
        if goal_match:
            minutes, seconds = _parse_time(goal_match.group("time"))
            scorer = goal_match.group("scorer").strip(" -—:")
            if not scorer:
                continue
            result.goals.append(
                Goal(
                    time_raw=goal_match.group("time"),
                    scorer=scorer,
                    half=current_half,
                    minutes=minutes,
                    seconds=seconds,
                )
            )

    return result
