import pandas as pd
from typing import Any, Optional


def format_lap_time(lap_time: Any) -> Optional[str]:
    """Formats a Timedelta into a M:SS.mmm string."""
    if pd.isna(lap_time):
        return None

    total_seconds = lap_time.total_seconds()
    minutes = int(total_seconds // 60)
    seconds = int(total_seconds % 60)
    milliseconds = int((total_seconds % 1) * 1000)

    return f"{minutes}:{seconds:02d}.{milliseconds:03d}"


def to_milliseconds(delta: Any) -> Optional[int]:
    """Converts a pandas Timedelta to total milliseconds."""
    if pd.isna(delta) or not hasattr(delta, 'total_seconds'):
        return None
    return int(delta.total_seconds() * 1000)


def get_scalar_value(series: pd.Series, key: str) -> Any:
    """Safely extracts a scalar value from a Pandas Series, handling potential nested series."""
    value = series.get(key)
    if value is not None and hasattr(value, 'iloc'):
        # If it's a Series/DataFrame (happens with duplicate indices), take the first item
        return value.iloc[0] if not value.empty else None
    return value
