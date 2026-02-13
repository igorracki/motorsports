import fastf1
import os
import logging
from typing import List, Optional
from .provider import Provider
from ..models import (
    RaceWeekend, SessionResult
)
from ..utils.extractors import extract_race_weekend, extract_driver_result

logger = logging.getLogger(__name__)

class FastF1Provider(Provider):

    def __init__(self):
        cache_directory = os.path.join(os.getcwd(), '.cache')
        if not os.path.exists(cache_directory):
            os.makedirs(cache_directory)
        fastf1.Cache.enable_cache(cache_directory)

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        logger.info(f"Fetching event schedule for year {year}")
        try:
            schedule = fastf1.get_event_schedule(year)
            race_weekends = []
            for _, weekend in schedule.iterrows():
                race_weekend = extract_race_weekend(weekend)
                if race_weekend:
                    race_weekends.append(race_weekend)
            logger.info(f"Found {len(race_weekends)} race weekends for {year}")
            return race_weekends
        except Exception as e:
            logger.error(f"Error fetching weekend events for year {year}: {e}")
            return []

    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        logger.info(f"Fetching session results: {year} Round {round_number} ({session_type})")
        try:
            session = fastf1.get_session(year, round_number, session_type)
            session.load(laps=True, telemetry=False,
                         weather=False, messages=False)

            if session.results is None or session.results.empty:
                logger.warning(f"No results found for {year} Round {round_number} ({session_type})")
                return None

            results = []
            for _, row in session.results.iterrows():
                driver_result = extract_driver_result(row, session, session_type)
                results.append(driver_result)

            logger.info(f"Session loaded. Processing {len(results)} drivers.")
            return SessionResult(
                year=year,
                round=round_number,
                session_type=session_type,
                results=results
            )
        except Exception as e:
            logger.error(f"Error fetching session results for {year} round {round_number} {session_type}: {e}")
            return None
