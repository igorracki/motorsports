import pandas as pd
import logging
from datetime import datetime
from typing import Optional
from ..models import RaceWeekend, Session
from .converters import to_datetime, datetime_to_ms

logger = logging.getLogger(__name__)

def extract_race_weekend(weekend: pd.Series) -> Optional[RaceWeekend]:
    event_name = str(weekend.get('EventName', "Unknown Event"))
    logger.info(f"Extracting race weekend: {event_name}")
    sessions = []
    for session_index in range(1, 6):
        session_type = str(weekend.get(f"Session{session_index}", ""))
        local_time = to_datetime(weekend.get(f"Session{session_index}Date"))
        utc_time = to_datetime(weekend.get(f"Session{session_index}DateUtc"))

        if isinstance(local_time, datetime) and isinstance(utc_time, datetime):
            utc_ms = datetime_to_ms(utc_time)
            
            # Calculate track offset in ms: (Local Time - UTC Time)
            # local_time from FastF1 is timezone-aware if available
            offset = local_time.utcoffset()
            offset_ms = int(offset.total_seconds() * 1000) if offset else 0
            
            if utc_ms is not None:
                sessions.append(Session(
                    type=session_type,
                    time_utc_ms=utc_ms,
                    utc_offset_ms=offset_ms
                ))

    raw_round = weekend.get('RoundNumber')
    round_number = int(raw_round) if raw_round is not None and pd.notna(
        raw_round) else 0

    race_weekend = RaceWeekend(
        round=round_number,
        full_name=str(weekend.get('OfficialEventName', "") or ""),
        name=str(weekend.get('EventName', "") or ""),
        location=str(weekend.get('Location', "") or ""),
        country=str(weekend.get('Country', "") or ""),
        sessions=sessions
    )
    logger.info(f"Successfully extracted race weekend: {race_weekend.name} with {len(sessions)} sessions")
    return race_weekend
