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
from ..utils.session_extractors import extract_race_weekend
from ..utils.result_extractors import extract_driver_result
from ..utils.circuit_extractors import (
    extract_circuit_layout, extract_circuit_location, extract_circuit_metrics
)

from ..utils.converters import datetime_to_ms, to_datetime

logger = logging.getLogger(__name__)

class FastF1Provider(Provider):

    def __init__(self):
        cache_directory = os.path.join(os.getcwd(), '.cache')
        if not os.path.exists(cache_directory):
            os.makedirs(cache_directory)
        fastf1.Cache.enable_cache(cache_directory)

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        logger.info(f"Entry: get_weekend_events(year={year})")
        try:
            schedule = fastf1.get_event_schedule(year)
            race_weekends = []
            for _, weekend in schedule.iterrows():
                race_weekend = extract_race_weekend(weekend)
                if race_weekend:
                    race_weekends.append(race_weekend)
            logger.info(f"Exit: get_weekend_events(year={year}) - Found {len(race_weekends)} race weekends")
            return race_weekends
        except Exception:
            logger.exception(f"Error fetching weekend events for year {year}")
            return []

    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        logger.info(f"Entry: get_session_results(year={year}, round={round_number}, session={session_type})")
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

            logger.info(f"Exit: get_session_results(year={year}, round={round_number}, session={session_type}) - Found {len(results)} drivers")
            return SessionResult(
                year=year,
                round=round_number,
                session_type=session_type,
                results=results
            )
        except Exception:
            logger.exception(f"Error fetching session results for {year} round {round_number} {session_type}")
            return None

    def get_circuit_data(self, year: int, round_number: int) -> Optional[Circuit]:
        logger.info(f"Entry: get_circuit_data(year={year}, round={round_number})")
        try:
            ergast = Ergast()
            ergast_circuits = ergast.get_circuits(season=year, round=round_number)
            location_info = extract_circuit_location(ergast_circuits)

            session = fastf1.get_session(year, round_number, 'Q')
            session.load(telemetry=True, laps=True, weather=False, messages=False)
            
            metrics = extract_circuit_metrics(session)
            
            rotation = 0.0
            try:
                circuit_info = session.get_circuit_info()
                if circuit_info is not None:
                    rotation = circuit_info.rotation
            except Exception:
                pass 

            layout = extract_circuit_layout(session)
            
            event_date = to_datetime(session.event.EventDate)
            event_date_ms = datetime_to_ms(event_date)
            if event_date_ms is None:
                event_date_ms = 0

            circuit = Circuit(
                circuit_name=location_info["circuit_name"],
                location=location_info["location"],
                country=location_info["country"],
                latitude=location_info["latitude"],
                longitude=location_info["longitude"],
                length_km=metrics["length_km"],
                corners=metrics["corners"],
                layout=layout,
                event_name=str(session.event.EventName),
                event_date_ms=event_date_ms,
                rotation=rotation
            )
            logger.info(f"Exit: get_circuit_data(year={year}, round={round_number}) - Success for {circuit.circuit_name}")
            return circuit

        except Exception:
            logger.exception(f"Error fetching circuit data for {year} Round {round_number}")
            return None
