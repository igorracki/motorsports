import json
from dataclasses import dataclass, asdict
from datetime import datetime
from typing import List, Optional
from .session import Session


@dataclass
class RaceWeekend:
    round: int
    full_name: str
    name: str
    location: str
    country: str
    sessions: List[Session]
    event_format: str

    def json(self) -> str:
        return json.dumps(asdict(self), indent=2)
