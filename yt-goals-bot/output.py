"""Сборка результата: картинка + подпись с YouTube-ссылками."""

from __future__ import annotations

from dataclasses import dataclass

from generator import generate_comment
from image_generator import generate_match_image
from parser import MatchResult


@dataclass
class MatchOutput:
    image_png: bytes
    caption: str


def build_match_output(
    match: MatchResult,
    first_half_url: str,
    second_half_url: str,
) -> MatchOutput:
    image_png = generate_match_image(match)
    caption = generate_comment(match, first_half_url, second_half_url)
    if len(caption) > 1024:
        caption = caption[:1020] + "..."
    return MatchOutput(image_png=image_png, caption=caption)
