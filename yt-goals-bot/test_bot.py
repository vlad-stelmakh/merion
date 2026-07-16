"""Тесты парсера и генератора."""

import unittest

from generator import generate_comment, minute_to_seconds, seconds_to_timestamp
from image_generator import generate_match_image
from parser import parse_match_results


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
        self.assertEqual(len(first), 3)
        self.assertEqual(len(second), 6)

    def test_team_sides_pakuta(self) -> None:
        match = parse_match_results(PAKUTA_TEXT)
        home = [g for g in match.goals if g.side == "home"]
        away = [g for g in match.goals if g.side == "away"]
        self.assertEqual(len(home), 6)
        self.assertEqual(len(away), 2)


class GeneratorTests(unittest.TestCase):
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
        self.assertIn("AAAA1111", comment)
        self.assertIn("BBBB2222", comment)
        self.assertIn("14:50", comment)  # 15' - 10 сек
        self.assertIn("Filipp PlusMinus (Sandro Karchava)", comment)

    def test_generate_match_image(self) -> None:
        match = parse_match_results(PAKUTA_TEXT)
        png = generate_match_image(match)
        self.assertTrue(png.startswith(b"\x89PNG"))


if __name__ == "__main__":
    unittest.main()
