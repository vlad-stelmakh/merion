"""Генерация PNG-карточки с результатом матча."""

from __future__ import annotations

import io
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

from parser import Goal, MatchResult

# Цвета в стиле примера
BG_COLOR = (34, 120, 68)
HEADER_BG = (20, 20, 20)
HEADER_TEXT = (255, 255, 255)
SCORE_COLOR = (255, 255, 255)
TEAM_COLOR = (255, 255, 255)
GOAL_COLOR = (255, 255, 255)
LOGO_BORDER = (255, 255, 255)
LOGO_FILL_LEFT = (180, 40, 40)
LOGO_FILL_RIGHT = (160, 60, 180)

LEAGUE_TITLE = "F.F.F."
LEAGUE_SUBTITLE = "FRIENDS FOR FRIENDS"

WIDTH = 960
PADDING = 36
GOAL_LINE_HEIGHT = 34
HEADER_HEIGHT = 72
TEAMS_ROW_HEIGHT = 150
GOAL_BALL = ""  # рисуем мяч отдельно

FONT_DIR = Path("/usr/share/fonts/truetype/dejavu")


def _font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont | ImageFont.ImageFont:
    names = (
        ["DejaVuSans-Bold.ttf", "DejaVuSans.ttf"]
        if bold
        else ["DejaVuSans.ttf", "DejaVuSans-Bold.ttf"]
    )
    for name in names:
        path = FONT_DIR / name
        if path.exists():
            return ImageFont.truetype(str(path), size=size)
    return ImageFont.load_default()


def _split_scorer_assist(scorer: str) -> tuple[str, str | None]:
    scorer = scorer.strip()
    if "(" in scorer and scorer.endswith(")"):
        name, rest = scorer.rsplit("(", 1)
        return name.strip(), f"({rest}"
    return scorer, None


def _goals_for_side(match: MatchResult, side: str) -> list[Goal]:
    tagged = [g for g in match.goals if g.side == side]
    if tagged:
        return sorted(tagged, key=lambda g: g.minute)

    if side == "home" and match.score_home is not None:
        ordered = sorted(match.goals, key=lambda g: g.minute)
        return ordered[: match.score_home]
    if side == "away" and match.score_away is not None:
        ordered = sorted(match.goals, key=lambda g: g.minute)
        return ordered[-match.score_away :]
    return []


def _draw_logo(
    draw: ImageDraw.ImageDraw,
    center_x: int,
    center_y: int,
    radius: int,
    fill: tuple[int, int, int],
    text: str,
) -> None:
    x0, y0 = center_x - radius, center_y - radius
    x1, y1 = center_x + radius, center_y + radius
    draw.ellipse((x0, y0, x1, y1), fill=fill, outline=LOGO_BORDER, width=3)
    font = _font(max(18, radius // 2), bold=True)
    bbox = draw.textbbox((0, 0), text, font=font)
    tw = bbox[2] - bbox[0]
    th = bbox[3] - bbox[1]
    draw.text((center_x - tw // 2, center_y - th // 2 - 2), text, fill=HEADER_TEXT, font=font)


def _draw_ball(draw: ImageDraw.ImageDraw, x: int, y: int, size: int = 10) -> int:
    draw.ellipse((x, y, x + size, y + size), fill=(255, 255, 255), outline=(30, 30, 30), width=1)
    return x + size + 6


def _fit_team_name(
    draw: ImageDraw.ImageDraw,
    name: str,
    max_width: int,
    base_size: int = 24,
) -> tuple[ImageFont.ImageFont, str]:
    for size in range(base_size, 14, -2):
        font = _font(size, bold=True)
        if _text_width(draw, name, font) <= max_width:
            return font, name
def _team_initials(name: str) -> str:
    parts = [p for p in name.replace(".", " ").split() if p]
    if not parts:
        return "?"
    if len(parts) == 1:
        return parts[0][:2].upper()
    return (parts[0][0] + parts[1][0]).upper()


def _text_width(draw: ImageDraw.ImageDraw, text: str, font: ImageFont.ImageFont) -> int:
    bbox = draw.textbbox((0, 0), text, font=font)
    return bbox[2] - bbox[0]


def _goal_line_parts(minute: int, scorer: str) -> tuple[str, str | None]:
    name, assist = _split_scorer_assist(scorer)
    return f"{minute}'", name + (f" {assist}" if assist else "")


def _draw_goal_line(
    draw: ImageDraw.ImageDraw,
    x: int,
    y: int,
    minute: int,
    scorer: str,
    font: ImageFont.ImageFont,
    max_width: int,
) -> int:
    minute_text, player_text = _goal_line_parts(minute, scorer)
    draw.text((x, y), minute_text, fill=GOAL_COLOR, font=font)
    cursor = x + _text_width(draw, minute_text + " ", font) + 2
    cursor = _draw_ball(draw, cursor, y + 6)
    draw.text((cursor, y), player_text, fill=GOAL_COLOR, font=font)
    return y + GOAL_LINE_HEIGHT


def generate_match_image(match: MatchResult, league_title: str = LEAGUE_TITLE) -> bytes:
    """Рисует карточку матча и возвращает PNG в bytes."""
    if not match.goals:
        raise ValueError("Нет голов для карточки.")

    home_team = match.home_team or "Команда 1"
    away_team = match.away_team or "Команда 2"
    score_home = match.score_home if match.score_home is not None else len(_goals_for_side(match, "home"))
    score_away = match.score_away if match.score_away is not None else len(_goals_for_side(match, "away"))

    home_goals = _goals_for_side(match, "home")
    away_goals = _goals_for_side(match, "away")
    max_goal_rows = max(len(home_goals), len(away_goals), 1)

    goal_font = _font(20)
    score_font = _font(56, bold=True)
    header_font = _font(22, bold=True)
    header_sub_font = _font(14)

    col_width = (WIDTH - PADDING * 2 - 80) // 2

    height = (
        PADDING
        + HEADER_HEIGHT
        + TEAMS_ROW_HEIGHT
        + max_goal_rows * GOAL_LINE_HEIGHT
        + PADDING
        + 20
    )

    img = Image.new("RGB", (WIDTH, height), BG_COLOR)
    draw = ImageDraw.Draw(img)

    # Шапка лиги
    header_w = 320
    header_h = 56
    header_x = (WIDTH - header_w) // 2
    header_y = PADDING
    draw.rounded_rectangle(
        (header_x, header_y, header_x + header_w, header_y + header_h),
        radius=8,
        fill=HEADER_BG,
    )
    title_bbox = draw.textbbox((0, 0), league_title, font=header_font)
    title_w = title_bbox[2] - title_bbox[0]
    draw.text(
        (WIDTH // 2 - title_w // 2, header_y + 8),
        league_title,
        fill=HEADER_TEXT,
        font=header_font,
    )
    sub_bbox = draw.textbbox((0, 0), LEAGUE_SUBTITLE, font=header_sub_font)
    sub_w = sub_bbox[2] - sub_bbox[0]
    draw.text(
        (WIDTH // 2 - sub_w // 2, header_y + 32),
        LEAGUE_SUBTITLE,
        fill=HEADER_TEXT,
        font=header_sub_font,
    )

    teams_y = header_y + HEADER_HEIGHT + 10
    left_x = PADDING + 70
    right_x = WIDTH - PADDING - 70
    center_x = WIDTH // 2

    _draw_logo(draw, left_x, teams_y + 35, 48, LOGO_FILL_LEFT, _team_initials(home_team))
    _draw_logo(draw, right_x, teams_y + 35, 48, LOGO_FILL_RIGHT, _team_initials(away_team))

    # Названия команд
    name_max_w = (WIDTH // 2) - 120
    for team, x in ((home_team, left_x), (away_team, right_x)):
        team_font, fitted = _fit_team_name(draw, team, name_max_w)
        tw = _text_width(draw, fitted, team_font)
        draw.text((x - tw // 2, teams_y + 95), fitted, fill=TEAM_COLOR, font=team_font)

    score_text = f"{score_home}  -  {score_away}"
    score_bbox = draw.textbbox((0, 0), score_text, font=score_font)
    score_w = score_bbox[2] - score_bbox[0]
    draw.text(
        (center_x - score_w // 2, teams_y + 30),
        score_text,
        fill=SCORE_COLOR,
        font=score_font,
    )

    goals_top = teams_y + TEAMS_ROW_HEIGHT
    left_col_x = PADDING
    right_col_x = WIDTH // 2 + 20

    def draw_goal_column(goals: list[Goal], x: int) -> None:
        y = goals_top
        for goal in goals:
            y = _draw_goal_line(draw, x, y, goal.minute, goal.scorer, goal_font, col_width)

    draw_goal_column(home_goals, left_col_x)
    draw_goal_column(away_goals, right_col_x)

    # Вертикальный разделитель
    draw.line(
        (WIDTH // 2, goals_top - 8, WIDTH // 2, height - PADDING),
        fill=(255, 255, 255, 80),
        width=2,
    )

    buffer = io.BytesIO()
    img.save(buffer, format="PNG", optimize=True)
    return buffer.getvalue()
