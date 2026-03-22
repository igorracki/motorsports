
import unittest
from unittest.mock import MagicMock
import pandas as pd
from datetime import timedelta
from src.core.utils.result_extractors import extract_driver_result
from src.core.utils.session_extractors import extract_race_weekend
from src.core.utils.circuit_extractors import extract_circuit_location, extract_circuit_metrics, extract_circuit_layout
from src.core.models import DriverResult, RaceWeekend
from src.core.models.circuit import CircuitLayoutPoint

class TestExtractors(unittest.TestCase):

    def setUp(self):
        self.mock_session = MagicMock()
        self.mock_laps = MagicMock()
        self.mock_session.laps = self.mock_laps
        
    def test_extract_driver_result_race_winner(self):
        row = pd.Series({
            'Abbreviation': 'VER',
            'DriverNumber': '1',
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
        self.assertEqual(result.driver.id, 'VER')
        self.assertEqual(result.driver.number, '1')
        self.assertEqual(result.position, 1)
        self.assertEqual(result.gap_ms, 0)
        self.assertEqual(result.total_time_ms, 5400000)
        self.assertEqual(result.fastest_lap_ms, 80000)
        if result.race_details:
            self.assertEqual(result.race_details.grid_position, 1)

    def test_extract_driver_result_race_second_place(self):
        row = pd.Series({
            'Abbreviation': 'HAM',
            'DriverNumber': '44',
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

        self.assertEqual(result.driver.id, 'HAM')
        self.assertEqual(result.driver.number, '44')
        self.assertEqual(result.position, 2)
        self.assertEqual(result.gap_ms, 10000)
        self.assertIsNone(result.total_time_ms)
        
        if result.race_details:
            self.assertEqual(result.race_details.grid_position, 3)

    def test_extract_driver_result_qualifying(self):
        row = pd.Series({
            'Abbreviation': 'LEC',
            'DriverNumber': '16',
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

        self.assertEqual(result.driver.id, 'LEC')
        self.assertEqual(result.driver.number, '16')
        self.assertEqual(result.position, 1)
        self.assertIsNotNone(result.qualifying_details)
        if result.qualifying_details:
            self.assertEqual(result.qualifying_details.q1_ms, 90000)
            self.assertEqual(result.qualifying_details.q3_ms, 88000)

    def test_extract_circuit_location(self):
        df = pd.DataFrame([{
            'circuitName': 'Monaco',
            'locality': 'Monte Carlo',
            'country': 'Monaco',
            'lat': 43.7347,
            'long': 7.42056
        }])
        result = extract_circuit_location(df)
        self.assertEqual(result['circuit_name'], 'Monaco')
        self.assertEqual(result['latitude'], 43.7347)
        self.assertEqual(result['longitude'], 7.42056)

    def test_extract_circuit_metrics(self):
        mock_session = MagicMock()
        mock_ci = MagicMock()
        mock_ci.corners = [1, 2, 3]
        mock_session.get_circuit_info.return_value = mock_ci
        
        mock_lap = MagicMock()
        mock_telemetry = pd.DataFrame({
            'Distance': [0, 1000, 2000],
            'Speed': [100.0, 200.0, 300.0],
            'Z': [10.0, 20.0, 30.0]
        })
        mock_lap.get_telemetry.return_value = mock_telemetry
        mock_session.laps.pick_fastest.return_value = mock_lap
        
        result = extract_circuit_metrics(mock_session)
        self.assertEqual(result['corners'], 3)
        self.assertEqual(result['length_km'], 2.0)
        self.assertEqual(result['max_speed_kmh'], 300.0)
        self.assertEqual(result['max_altitude_m'], 30.0)
        self.assertEqual(result['min_altitude_m'], 10.0)

    def test_extract_circuit_layout(self):
        mock_session = MagicMock()
        mock_lap = MagicMock()
        # Mock telemetry with X, Y coordinates
        mock_telemetry = pd.DataFrame({
            'X': [100.0, 200.0],
            'Y': [300.0, 400.0]
        })
        mock_lap.get_telemetry.return_value = mock_telemetry
        mock_session.laps.pick_fastest.return_value = mock_lap
        
        result = extract_circuit_layout(mock_session)
        self.assertEqual(len(result), 2)
        self.assertEqual(result[0].x, 100.0)
        self.assertEqual(result[0].y, 300.0)

    def test_extract_race_weekend(self):
        # Australia is UTC+11 in March
        local_time = pd.Timestamp('2026-03-06 13:00:00', tz='Australia/Melbourne')
        utc_time = pd.Timestamp('2026-03-06 02:00:00')
        
        weekend_data = pd.Series({
            'RoundNumber': 1.0,
            'OfficialEventName': 'Australian Grand Prix',
            'EventName': 'Australian GP',
            'Location': 'Melbourne',
            'Country': 'Australia',
            'EventDate': pd.Timestamp('2026-03-08'),
            'Session1': 'Practice 1',
            'Session1Date': local_time,
            'Session1DateUtc': utc_time,
        })
        
        result = extract_race_weekend(weekend_data)
        self.assertIsInstance(result, RaceWeekend)
        self.assertEqual(result.round, 1)
        self.assertEqual(result.name, 'Australian GP')
        self.assertEqual(len(result.sessions), 1)
        self.assertEqual(result.sessions[0].type, 'Practice 1')
        
        # UTC 2026-03-06 02:00:00 is 1772762400000 ms
        self.assertEqual(result.sessions[0].time_utc_ms, 1772762400000)
        # Offset for Australia/Melbourne is +11h = 39600000 ms
        self.assertEqual(result.sessions[0].utc_offset_ms, 39600000)


if __name__ == '__main__':
    unittest.main()
