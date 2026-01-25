from dataclasses import dataclass
from datetime import datetime


@dataclass
class SessionInfo:
    name: str
    session_type: str
    date: datetime
    event_name: str
    location: str
    country: str


@dataclass
class Session:
    type: str
    time_local: datetime
    time_utc: datetime
