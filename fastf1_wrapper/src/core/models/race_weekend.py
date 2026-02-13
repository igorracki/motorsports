import json
from dataclasses import dataclass, asdict
from datetime import datetime
from typing import List
from .session import Session


@dataclass
class RaceWeekend:
    round: int
    full_name: str
    name: str
    location: str
    country: str
    start_date: datetime
    sessions: List[Session]

    def json(self) -> str:
        def convert_datetime(original):
            if isinstance(original, datetime):
                return original.isoformat()
            return original

        return json.dumps(asdict(self), default=convert_datetime, indent=2)
