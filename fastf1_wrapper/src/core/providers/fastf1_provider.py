import fastf1
from fastf1.ergast import Ergast
import os
import logging
from typing import List, Optional
from .provider import Provider
from ..models import (
    RaceWeekend, SessionResult
)
from ..models.circuit import Circuit, CircuitLayoutPoint
from ..utils.extractors import (
    extract_race_weekend, extract_driver_result, 
    extract_circuit_layout, extract_circuit_location, extract_circuit_metrics
)

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

    def get_circuit_data(self, year: int, round_number: int) -> Optional[Circuit]:
        logger.info(f"Fetching circuit data for {year} Round {round_number}")
        try:
            # 1. Fetch Location Data from Ergast
            ergast = Ergast()
            ergast_circuits = ergast.get_circuits(season=year, round=round_number)
            location_info = extract_circuit_location(ergast_circuits)

            # 2. Fetch Session Data (Qualifying for best layout)
            logger.info("Loading Qualifying session for layout extraction...")
            session = fastf1.get_session(year, round_number, 'Q')
            session.load(telemetry=True, laps=True, weather=False, messages=False)
            
            # 3. Extract Circuit Metrics (Corners, Length)
            metrics = extract_circuit_metrics(session)
            
            # 4. Extract Layout with Rotation
            rotation = 0.0
            try:
                circuit_info = session.get_circuit_info()
                if circuit_info is not None:
                    rotation = circuit_info.rotation
            except Exception:
                pass # Ignore if circuit info not available yet

            layout = extract_circuit_layout(session, rotation=rotation)
            
            return Circuit(
                circuit_name=location_info["circuit_name"],
                location=location_info["location"],
                country=location_info["country"],
                latitude=location_info["latitude"],
                longitude=location_info["longitude"],
                length_km=metrics["length_km"],
                corners=metrics["corners"],
                layout=layout,
                event_name=str(session.event.EventName),
                event_date=str(session.event.EventDate)
            )

        except Exception as e:
            logger.error(f"Error fetching circuit data: {e}")
            return None
