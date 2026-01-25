import json
from dataclasses import dataclass, asdict
from datetime import datetime
from typing import List
from .session import Session


@dataclass
class Event:
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


# RoundNumber                                                         24
# Country                                           United Arab Emirates
# Location                                                    Yas Island
# OfficialEventName    FORMULA 1 ETIHAD AIRWAYS ABU DHABI GRAND PRIX ...
# EventDate                                          2025-12-07 00:00:00
# EventName                                         Abu Dhabi Grand Prix
# EventFormat                                               conventional
# Session1                                                    Practice 1
# Session1Date                                 2025-12-05 13:30:00+04:00
# Session1DateUtc                                    2025-12-05 09:30:00
# Session2                                                    Practice 2
# Session2Date                                 2025-12-05 17:00:00+04:00
# Session2DateUtc                                    2025-12-05 13:00:00
# Session3                                                    Practice 3
# Session3Date                                 2025-12-06 14:30:00+04:00
# Session3DateUtc                                    2025-12-06 10:30:00
# Session4                                                    Qualifying
# Session4Date                                 2025-12-06 18:00:00+04:00
# Session4DateUtc                                    2025-12-06 14:00:00
# Session5                                                          Race
# Session5Date                                 2025-12-07 17:00:00+04:00
# Session5DateUtc                                    2025-12-07 13:00:00
# F1ApiSupport                                                      True
