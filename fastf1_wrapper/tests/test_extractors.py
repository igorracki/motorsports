
import unittest
from unittest.mock import MagicMock
import pandas as pd
from datetime import timedelta
from src.core.utils.extractors import extract_driver_result
from src.core.models import DriverResult

class TestExtractors(unittest.TestCase):

    def setUp(self):
        self.mock_session = MagicMock()
        self.mock_laps = MagicMock()
        self.mock_session.laps = self.mock_laps
        
    def test_extract_driver_result_race_winner(self):
        row = pd.Series({
            'Abbreviation': 'VER',
            'FullName': 'Max Verstappen',
            'CountryCode': 'NED',
            'TeamName': 'Red Bull Racing',
            'Status': 'Finished',
            'Laps': 50,
            'Position': 1.0,
            'Time': pd.Timedelta(hours=1, minutes=30),
            'GridPosition': 1.0
        })

        mock_driver_laps = MagicMock()
        mock_driver_laps.empty = False
        self.mock_laps.pick_drivers.return_value = mock_driver_laps
        mock_driver_laps.pick_fastest.return_value = pd.Series({'LapTime': pd.Timedelta(minutes=1, seconds=20)})

        result = extract_driver_result(row, self.mock_session, 'R')

        self.assertIsInstance(result, DriverResult)
        self.assertEqual(result.position, 1)
        self.assertEqual(result.gap_ms, 0)
        self.assertEqual(result.total_time_ms, 5400000)
        self.assertEqual(result.fastest_lap_ms, 80000)
        if result.race_details:
            self.assertEqual(result.race_details.grid_position, 1)

    def test_extract_driver_result_race_second_place(self):
        row = pd.Series({
            'Abbreviation': 'HAM',
            'FullName': 'Lewis Hamilton',
            'CountryCode': 'GBR',
            'TeamName': 'Mercedes',
            'Status': 'Finished',
            'Laps': 50,
            'Position': 2.0,
            'Time': pd.Timedelta(seconds=10),
            'GridPosition': 3.0
        })
        
        mock_driver_laps = MagicMock()
        mock_driver_laps.empty = False
        self.mock_laps.pick_drivers.return_value = mock_driver_laps
        mock_driver_laps.pick_fastest.return_value = pd.Series({'LapTime': pd.Timedelta(minutes=1, seconds=21)})


        result = extract_driver_result(row, self.mock_session, 'R')

        self.assertEqual(result.position, 2)
        self.assertEqual(result.gap_ms, 10000)
        self.assertIsNone(result.total_time_ms)
        
        if result.race_details:
            self.assertEqual(result.race_details.grid_position, 3)

    def test_extract_driver_result_qualifying(self):
        row = pd.Series({
            'Abbreviation': 'LEC',
            'FullName': 'Charles Leclerc',
            'CountryCode': 'MON',
            'TeamName': 'Ferrari',
            'Status': 'Finished',
            'Laps': 15,
            'Position': 1.0,
            'Q1': pd.Timedelta(minutes=1, seconds=30),
            'Q2': pd.Timedelta(minutes=1, seconds=29),
            'Q3': pd.Timedelta(minutes=1, seconds=28),
        })

        result = extract_driver_result(row, self.mock_session, 'Q')

        self.assertEqual(result.position, 1)
        self.assertIsNotNone(result.qualifying_details)
        if result.qualifying_details:
            self.assertEqual(result.qualifying_details.q1_ms, 90000)
            self.assertEqual(result.qualifying_details.q3_ms, 88000)


if __name__ == '__main__':
    unittest.main()
