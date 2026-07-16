"""Работа с YouTube-ссылками и таймкодами."""

from __future__ import annotations

import re
from urllib.parse import parse_qs, urlparse


YOUTUBE_RE = re.compile(
    r"(?i)(?:https?://)?(?:www\.)?"
    r"(?:youtube\.com/(?:watch\?v=|shorts/|live/|embed/)|youtu\.be/)"
    r"([A-Za-z0-9_-]{11})"
)


def extract_video_id(url: str) -> str | None:
    """Достаёт video id из youtube/youtu.be ссылки."""
    text = url.strip()
    match = YOUTUBE_RE.search(text)
    if match:
        return match.group(1)

    parsed = urlparse(text if "://" in text else f"https://{text}")
    host = (parsed.netloc or "").lower().removeprefix("www.")
    if host == "youtu.be":
        candidate = parsed.path.strip("/").split("/")[0]
        return candidate if len(candidate) == 11 else None
    if host in {"youtube.com", "m.youtube.com", "music.youtube.com"}:
        qs = parse_qs(parsed.query)
        if "v" in qs and len(qs["v"][0]) == 11:
            return qs["v"][0]
    return None


def is_youtube_url(text: str) -> bool:
    return extract_video_id(text) is not None


def timestamp_url(video_id: str, seconds: int) -> str:
    """Ссылка на момент ролика (секунды от начала)."""
    seconds = max(0, int(seconds))
    return f"https://youtu.be/{video_id}?t={seconds}"


def format_yt_time(seconds: int) -> str:
    """Секунды -> m:ss или h:mm:ss для отображения."""
    seconds = max(0, int(seconds))
    hours, rem = divmod(seconds, 3600)
    minutes, secs = divmod(rem, 60)
    if hours:
        return f"{hours}:{minutes:02d}:{secs:02d}"
    return f"{minutes}:{secs:02d}"
