"""Сборка результата для бота."""

from __future__ import annotations

from dataclasses import dataclass

from generator import generate_comment
from image_generator import generate_match_image
from parser import MatchResult


@dataclass
class MatchOutput:
    caption: str
    image_png: bytes | None = None


def build_comment_output(
    match: MatchResult,
    first_half_url: str,
    second_half_url: str,
) -> MatchOutput:
    caption = generate_comment(match, first_half_url, second_half_url)
    return MatchOutput(caption=caption)


def build_match_output(
    match: MatchResult,
    first_half_url: str,
    second_half_url: str,
) -> MatchOutput:
    caption = generate_comment(match, first_half_url, second_half_url)
    image_png = generate_match_image(match)
    return MatchOutput(caption=caption, image_png=image_png)
