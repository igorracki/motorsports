import json
from datetime import datetime
from dataclasses import dataclass, asdict
from .driver import Driver
from .session import SessionInfo
from typing import List


@dataclass
class DriverResult:
    position: int
    driver: Driver
    fastest_lap_time: str
    positions_gained: int


@dataclass
class RaceResult:
    session_info: SessionInfo
    results: List[DriverResult]

    def json(self) -> str:
        def convert_datetime(original):
            if isinstance(original, datetime):
                return original.isoformat()
            return original

        return json.dumps(asdict(self), default=convert_datetime, indent=2)
