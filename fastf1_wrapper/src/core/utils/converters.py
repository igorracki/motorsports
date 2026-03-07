import pandas as pd
from typing import Any, Optional
from datetime import datetime, timezone


def get_scalar_value(series: pd.Series, key: str) -> Any:
    value = series.get(key)
    if value is not None and hasattr(value, 'iloc'):
        return value.iloc[0] if not value.empty else None
    return value


def to_milliseconds(delta: Any) -> Optional[int]:
    if pd.isna(delta):
        return None
        
    try:
        # If it's already a timedelta-like object with total_seconds
        if hasattr(delta, 'total_seconds'):
            return int(delta.total_seconds() * 1000)
            
        # Try to convert to timedelta
        td = pd.to_timedelta(delta)
        if pd.isna(td):
            return None
        return int(td.total_seconds() * 1000)
    except (ValueError, TypeError, AttributeError):
        return None


def to_datetime(value: Any) -> Optional[datetime]:
    if value is None:
        return None
    
    # Handle pandas NaT
    if pd.isna(value):
        return None
        
    if hasattr(value, 'to_pydatetime'):
        dt = value.to_pydatetime()
        if dt is None or pd.isna(dt):
            return None
        return dt
        
    if isinstance(value, datetime):
        return value
    return None


def datetime_to_ms(date_time: Optional[datetime]) -> Optional[int]:
    if date_time is None:
        return None
    
    try:
        # pd.NaT is an instance of datetime but raises ValueError on timestamp()
        return int(date_time.timestamp() * 1000)
    except (ValueError, AttributeError):
        return None

