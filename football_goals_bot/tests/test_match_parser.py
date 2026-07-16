import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from match_parser import parse_results


SAMPLE = """
FC Kucha 3:6 FC Serega United

1 тайм:
15' Filipp PlusMinus (Sandro Karchava)
22' Filipp PlusMinus (Dima Semin)
40' Bobeeo — Sandro Karchava (автогол)

2 тайм:
11' Evgensky (Egor Levin)
20' Egor Levin
31' Evgensky
33' Andrei Tereshkov (Vlados Vlad)
40' Evgensky (Vlados Vlad)
49' Egor Levin (Levan Kviki)
"""


def test_parse_teams_and_score():
    match = parse_results(SAMPLE)
    assert match.home_team == "FC Kucha"
    assert match.away_team == "FC Serega United"
    assert match.score == "3:6"


def test_parse_goals_by_half():
    match = parse_results(SAMPLE)
    assert len(match.goals) == 9
    half1 = [g for g in match.goals if g.half == 1]
    half2 = [g for g in match.goals if g.half == 2]
    assert len(half1) == 3
    assert len(half2) == 6
    assert half1[0].minutes == 15
    assert half1[0].scorer == "Filipp PlusMinus (Sandro Karchava)"
    assert half2[-1].minutes == 49
    assert half2[-1].scorer == "Egor Levin (Levan Kviki)"


def test_parse_time_with_seconds():
    text = "1 тайм:\n15:30 Player One\n2 тайм:\n1:05 Player Two"
    match = parse_results(text)
    assert match.goals[0].total_seconds == 15 * 60 + 30
    assert match.goals[1].half == 2
    assert match.goals[1].total_seconds == 65
