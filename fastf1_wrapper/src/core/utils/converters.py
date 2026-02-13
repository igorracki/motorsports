import pandas as pd
from typing import Any, Optional
from datetime import datetime


def get_scalar_value(series: pd.Series, key: str) -> Any:
    value = series.get(key)
    if value is not None and hasattr(value, 'iloc'):
        return value.iloc[0] if not value.empty else None
    return value


def to_milliseconds(delta: Any) -> Optional[int]:
    if pd.isna(delta) or not hasattr(delta, 'total_seconds'):
        return None
    return int(delta.total_seconds() * 1000)


def to_datetime(value: Any) -> Optional[datetime]:
    if value is not None and hasattr(value, 'to_pydatetime'):
        return value.to_pydatetime()
    if isinstance(value, datetime):
        return value
    return None
