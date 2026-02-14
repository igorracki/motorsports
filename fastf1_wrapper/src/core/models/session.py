from dataclasses import dataclass
from datetime import datetime


@dataclass
class Session:
    type: str
    time_local_ms: int
    time_utc_ms: int
