"""Генерация YouTube-комментария с таймкодами голов."""

from __future__ import annotations

import os
import re
from urllib.parse import parse_qs, urlparse

from parser import Goal, Half, MatchResult

OFFSET_SECONDS = int(os.getenv("YT_OFFSET_SECONDS", "10"))
HALF_MINUTES = int(os.getenv("HALF_MINUTES", "25"))


def goal_video_minute(goal: Goal, half_minutes: int = HALF_MINUTES) -> int:
    """Минута внутри конкретного YouTube-видео тайма."""
    if goal.half == Half.SECOND and goal.minute > half_minutes:
        return goal.minute - half_minutes
    return goal.minute


def goal_to_seconds(goal: Goal) -> int:
    video_minute = goal_video_minute(goal)
    return max(0, video_minute * 60 - OFFSET_SECONDS)


def minute_to_seconds(minute: int) -> int:
    return max(0, minute * 60 - OFFSET_SECONDS)


def seconds_to_timestamp(seconds: int) -> str:
    minutes = seconds // 60
    secs = seconds % 60
    if minutes >= 60:
        hours = minutes // 60
        minutes = minutes % 60
        return f"{hours}:{minutes:02d}:{secs:02d}"
    return f"{minutes}:{secs:02d}"


def goal_to_timestamp(goal: Goal) -> str:
    return seconds_to_timestamp(goal_to_seconds(goal))


def extract_video_id(url: str) -> str | None:
    parsed = urlparse(url.strip())
    host = parsed.netloc.lower()

    if "youtu.be" in host:
        video_id = parsed.path.lstrip("/").split("/")[0]
        return video_id or None

    if "youtube.com" in host:
        if parsed.path == "/watch":
            query = parse_qs(parsed.query)
            video_ids = query.get("v", [])
            return video_ids[0] if video_ids else None
        match = re.match(r"^/(?:embed|shorts|live)/([^/?]+)", parsed.path)
        if match:
            return match.group(1)

    return None


def build_youtube_link(video_id: str, seconds: int) -> str:
    return f"https://www.youtube.com/watch?v={video_id}&t={seconds}s"


def generate_comment(
    match: MatchResult,
    first_half_url: str,
    second_half_url: str,
) -> str:
    if not match.goals:
        raise ValueError("Не найдено ни одного гола. Проверьте формат текста с результатами.")

    first_id = extract_video_id(first_half_url)
    second_id = extract_video_id(second_half_url)
    if not first_id:
        raise ValueError("Некорректная ссылка на 1-й тайм.")
    if not second_id:
        raise ValueError("Некорректная ссылка на 2-й тайм.")

    if match.home_team and match.away_team:
        if match.score_home is not None and match.score_away is not None:
            header = (
                f"⚽ {match.home_team} {match.score_home}:{match.score_away} "
                f"{match.away_team}\n"
            )
        else:
            header = f"⚽ {match.home_team} — {match.away_team}\n"
    else:
        header = "⚽ Голы матча\n"

    lines = [header.rstrip(), ""]

    first_half_goals = [g for g in match.goals if g.half == Half.FIRST]
    second_half_goals = [g for g in match.goals if g.half == Half.SECOND]

    if first_half_goals:
        lines.append("1-й тайм")
        for goal in first_half_goals:
            lines.append(f"{goal_to_timestamp(goal)} — {goal.scorer}")
        lines.append("")

    if second_half_goals:
        lines.append("2-й тайм")
        for goal in second_half_goals:
            lines.append(f"{goal_to_timestamp(goal)} — {goal.scorer}")

    return "\n".join(lines).strip()
