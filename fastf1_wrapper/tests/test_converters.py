
import unittest
import pandas as pd
from datetime import datetime, timedelta
from src.core.utils.converters import to_milliseconds, to_datetime, get_scalar_value

class TestConverters(unittest.TestCase):

    def test_to_milliseconds(self):
        delta = timedelta(seconds=1, milliseconds=500)
        self.assertEqual(to_milliseconds(delta), 1500)

        pd_delta = pd.Timedelta(seconds=2, milliseconds=100)
        self.assertEqual(to_milliseconds(pd_delta), 2100)

        self.assertIsNone(to_milliseconds(None))
        self.assertIsNone(to_milliseconds(pd.NaT))

    def test_to_datetime(self):
        dt = datetime(2023, 1, 1, 12, 0, 0)
        self.assertEqual(to_datetime(dt), dt)

        ts = pd.Timestamp("2023-01-01 12:00:00")
        self.assertEqual(to_datetime(ts), dt)

        self.assertIsNone(to_datetime(None))

    def test_get_scalar_value(self):
        series = pd.Series({'a': 1, 'b': 2})
        self.assertEqual(get_scalar_value(series, 'a'), 1)
        self.assertEqual(get_scalar_value(series, 'b'), 2)
        self.assertIsNone(get_scalar_value(series, 'c'))

if __name__ == '__main__':
    unittest.main()
