"""Тесты парсера и генератора."""

import unittest

from generator import generate_comment, goal_to_timestamp, goal_video_minute, minute_to_seconds, seconds_to_timestamp
from image_parser import _half_from_minute, _parse_goals_from_text
from image_generator import generate_match_image
from parser import Half, parse_match_results


EXAMPLE_TEXT = """FC Kucha 3:6 FC Serega United

15' Filipp PlusMinus (Sandro Karchava)
22' Filipp PlusMinus (Dima Semin)
40' Bobeeo — Sandro Karchava (автогол)

Evgensky 11', 20', 31'
Andrei Tereshkov 33'
Egor Levin 20', 49' (Levan Kviki)"""

PAKUTA_TEXT = """YFC Pakuta 6:2 City United FC

4' David Benidze (Vova Orange)
8' Vova Orange (David Benidze)
18' Luka Kaladze (Arthur Parkour)
23' Arthur Parkour (Luka Kaladze)
38' Arman (David Benidze)
45' Luka Kaladze (Arthur Parkour)

6' Efim Tarasenko (Danya Shangin)
20' Efim Tarasenko (Danya Shangin)"""


class ParserTests(unittest.TestCase):
    def test_parse_example_match(self) -> None:
        match = parse_match_results(EXAMPLE_TEXT)
        self.assertEqual(match.home_team, "FC Kucha")
        self.assertEqual(match.away_team, "FC Serega United")
        self.assertEqual(match.score_home, 3)
        self.assertEqual(match.score_away, 6)
        self.assertEqual(len(match.goals), 9)

    def test_first_half_goals_count(self) -> None:
        match = parse_match_results(EXAMPLE_TEXT)
        first = [g for g in match.goals if g.half.value == 1]
        second = [g for g in match.goals if g.half.value == 2]
        self.assertEqual(len(first), 2)  # 15', 22' (40' > 25 → 2-й тайм)
        self.assertEqual(len(second), 7)

    def test_team_sides_pakuta(self) -> None:
        match = parse_match_results(PAKUTA_TEXT)
        home = [g for g in match.goals if g.side == "home"]
        away = [g for g in match.goals if g.side == "away"]
        self.assertEqual(len(home), 6)
        self.assertEqual(len(away), 2)


HALF_FORMAT_TEXT = """1-й тайм

6' Олег Степанов (Дима Дидимер)
13' Серго (Дима Дидимер)
16' Женя Ответра
20' Женя Ответра
21' Женя Ответра (Олег Степанов)
30' Миша Юров (Сандро Карчава)
33' Олег Степанов (Женя Ответра)
35' Соута
40' Александр Косенков (Андрей Раб)

2-й тайм

41' Александр Косенков
48' Александр Косенков"""


class ImageParserTests(unittest.TestCase):
    def test_parse_goal_line(self) -> None:
        self.assertEqual(_parse_goals_from_text("6' Олег Степанов (Дима Дидимер)", "home")[0].minute, 6)

    def test_reject_noise(self) -> None:
        goals = _parse_goals_from_text("8' 2 KUCHA 41' Aleksandr Kosenkov", "away")
        self.assertEqual(len(goals), 1)
        self.assertEqual(goals[0].minute, 41)
        self.assertIn("Kosenkov", goals[0].scorer)
        self.assertEqual(_parse_goals_from_text("0' Aleksandr K senk v", "home"), [])

    def test_split_merged_goals(self) -> None:
        text = "40' Aleksandr Kosenkov (Andrei Rab) 41' Aleksandr Kosenkov"
        goals = _parse_goals_from_text(text, "home")
        self.assertEqual(len(goals), 2)
        self.assertEqual(goals[1].minute, 41)

    def test_extract_ocr_texts_formats(self) -> None:
        from image_parser import _extract_ocr_texts

        three = [([[0, 0], [1, 0], [1, 1], [0, 1]], "Goal", 0.9)]
        two = [[[[0, 0], [1, 0], [1, 1], [0, 1]], "Goal"]]
        self.assertEqual(_extract_ocr_texts(three), ["Goal"])
        self.assertEqual(_extract_ocr_texts(two), ["Goal"])

        self.assertEqual(_half_from_minute(21), Half.FIRST)
        self.assertEqual(_half_from_minute(41), Half.SECOND)


class HalfFormatTests(unittest.TestCase):
    def test_parse_half_sections(self) -> None:
        match = parse_match_results(HALF_FORMAT_TEXT)
        self.assertEqual(len(match.goals), 11)
        first = [g for g in match.goals if g.half == Half.FIRST]
        second = [g for g in match.goals if g.half == Half.SECOND]
        self.assertEqual(len(first), 9)
        self.assertEqual(len(second), 2)
        self.assertEqual(second[0].minute, 41)
        self.assertEqual(second[1].minute, 48)

    def test_half_format_youtube_timestamps(self) -> None:
        match = parse_match_results(HALF_FORMAT_TEXT)
        comment = generate_comment(match, "https://youtu.be/a", "https://youtu.be/b")
        self.assertIn("5:50 — Олег Степанов", comment)
        self.assertIn("39:50 — Александр Косенков (Андрей Раб)", comment)
        self.assertIn("15:50 — Александр Косенков", comment)
        self.assertIn("22:50 — Александр Косенков", comment)
        self.assertNotIn("40:50", comment)

    def test_second_half_video_minute(self) -> None:
        match = parse_match_results(HALF_FORMAT_TEXT)
        second = [g for g in match.goals if g.half == Half.SECOND]
        self.assertEqual(goal_video_minute(second[0]), 16)  # 41' → 16' видео
        self.assertEqual(goal_to_timestamp(second[0]), "15:50")


    def test_timestamp_offset(self) -> None:
        self.assertEqual(minute_to_seconds(15), 15 * 60 - 10)
        self.assertEqual(seconds_to_timestamp(650), "10:50")

    def test_generate_comment_contains_links(self) -> None:
        match = parse_match_results(EXAMPLE_TEXT)
        comment = generate_comment(
            match,
            "https://www.youtube.com/watch?v=AAAA1111",
            "https://www.youtube.com/watch?v=BBBB2222",
        )
        self.assertIn("FC Kucha 3:6 FC Serega United", comment)
        self.assertIn("14:50 — Filipp PlusMinus", comment)
        self.assertIn("10:50 — Evgensky", comment)
        self.assertNotIn("15'", comment)

    def test_generate_match_image(self) -> None:
        match = parse_match_results(PAKUTA_TEXT)
        png = generate_match_image(match)
        self.assertTrue(png.startswith(b"\x89PNG"))


if __name__ == "__main__":
    unittest.main()
