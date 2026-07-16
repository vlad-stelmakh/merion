import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from comment import generate_comment, goal_video_seconds
from match_parser import Goal, MatchResult, parse_results
from youtube import extract_video_id, format_yt_time, timestamp_url


def test_offset_ten_seconds():
    goal = Goal("15'", "Player", half=1, minutes=15, seconds=0)
    assert goal_video_seconds(goal, 10) == 14 * 60 + 50
    assert format_yt_time(goal_video_seconds(goal, 10)) == "14:50"


def test_offset_floor_zero():
    goal = Goal("0:05", "Player", half=1, minutes=0, seconds=5)
    assert goal_video_seconds(goal, 10) == 0


def test_youtube_id_and_url():
    assert extract_video_id("https://youtu.be/dQw4w9WgXcQ") == "dQw4w9WgXcQ"
    assert (
        extract_video_id("https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=12")
        == "dQw4w9WgXcQ"
    )
    assert timestamp_url("dQw4w9WgXcQ", 890) == "https://youtu.be/dQw4w9WgXcQ?t=890"


def test_generate_comment_links_by_half():
    text = """
    Home 1:1 Away
    1 тайм:
    15' Alpha
    2 тайм:
    11' Beta
    """
    match = parse_results(text)
    comment = generate_comment(match, "AAAAAAAAAAA", "BBBBBBBBBBB", offset_sec=10)

    assert "Home — Away" in comment
    assert "1 тайм:" in comment
    assert "2 тайм:" in comment
    assert "https://youtu.be/AAAAAAAAAAA?t=890" in comment  # 15:00 - 10
    assert "https://youtu.be/BBBBBBBBBBB?t=650" in comment  # 11:00 - 10
    assert "Alpha" in comment
    assert "Beta" in comment


def test_empty_match():
    match = MatchResult()
    comment = generate_comment(match, "AAAAAAAAAAA", "BBBBBBBBBBB")
    assert "Голы не найдены" in comment
