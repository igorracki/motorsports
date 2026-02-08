from dataclasses import dataclass
from datetime import datetime


@dataclass
class Session:
    type: str
    time_local: datetime
    time_utc: datetime
