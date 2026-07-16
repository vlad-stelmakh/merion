"""Генерация комментария с таймкодами голов (−offset секунд)."""

from __future__ import annotations

from match_parser import Goal, MatchResult
from youtube import format_yt_time, timestamp_url


DEFAULT_OFFSET_SEC = 10


def goal_video_seconds(goal: Goal, offset_sec: int = DEFAULT_OFFSET_SEC) -> int:
    """Время гола на видео минус offset (не ниже 0)."""
    return max(0, goal.total_seconds - offset_sec)


def generate_comment(
    match: MatchResult,
    half1_video_id: str,
    half2_video_id: str,
    offset_sec: int = DEFAULT_OFFSET_SEC,
) -> str:
    """
    Готовый текст для копирования в YouTube-комментарий.

    Для каждого гола — кликабельная ссылка на нужный тайм (−offset сек).
    """
    title = f"⚽ Голы матча {match.home_team} — {match.away_team}"
    if match.score:
        title += f" ({match.score})"

    lines = [title, ""]

    half1 = [g for g in match.goals if g.half == 1]
    half2 = [g for g in match.goals if g.half == 2]

    if half1:
        lines.append("1 тайм:")
        for goal in half1:
            lines.append(_goal_line(goal, half1_video_id, offset_sec))
        lines.append("")

    if half2:
        lines.append("2 тайм:")
        for goal in half2:
            lines.append(_goal_line(goal, half2_video_id, offset_sec))
        lines.append("")

    if not half1 and not half2:
        lines.append("Голы не найдены — проверь формат текста.")

    return "\n".join(lines).rstrip() + "\n"


def generate_telegram_html(
    match: MatchResult,
    half1_video_id: str,
    half2_video_id: str,
    offset_sec: int = DEFAULT_OFFSET_SEC,
) -> str:
    """HTML-версия для красивого сообщения в Telegram."""
    title = (
        f"⚽ <b>Голы матча { _esc(match.home_team) } — "
        f"{_esc(match.away_team)}</b>"
    )
    if match.score:
        title += f" ({_esc(match.score)})"

    blocks = [title, ""]

    half1 = [g for g in match.goals if g.half == 1]
    half2 = [g for g in match.goals if g.half == 2]

    if half1:
        blocks.append("<b>1 тайм</b>")
        for goal in half1:
            blocks.append(_goal_line_html(goal, half1_video_id, offset_sec))
        blocks.append("")

    if half2:
        blocks.append("<b>2 тайм</b>")
        for goal in half2:
            blocks.append(_goal_line_html(goal, half2_video_id, offset_sec))
        blocks.append("")

    if not half1 and not half2:
        blocks.append("Голы не найдены — проверь формат текста.")

    return "\n".join(blocks).rstrip()


def _goal_line(goal: Goal, video_id: str, offset_sec: int) -> str:
    sec = goal_video_seconds(goal, offset_sec)
    stamp = format_yt_time(sec)
    url = timestamp_url(video_id, sec)
    return f"{stamp} {url} — {goal.scorer}"


def _goal_line_html(goal: Goal, video_id: str, offset_sec: int) -> str:
    sec = goal_video_seconds(goal, offset_sec)
    stamp = format_yt_time(sec)
    url = timestamp_url(video_id, sec)
    return f'<a href="{url}">{stamp}</a> — {_esc(goal.scorer)}'


def _esc(text: str) -> str:
    return (
        text.replace("&", "&amp;")
        .replace("<", "&lt;")
        .replace(">", "&gt;")
    )
