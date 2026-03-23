import fastf1
from fastf1.ergast import Ergast
import pandas as pd
import os
import logging
import time
import unicodedata
from typing import Any, List, Optional
from .provider import Provider

from ..models import (
    RaceWeekend, SessionResult, DriverInfo
)
from ..models.circuit import Circuit, CircuitLayoutPoint
from ..utils.session_extractors import extract_race_weekend
from ..utils.result_extractors import extract_driver_result, extract_driver_info
from ..utils.circuit_extractors import (
    extract_circuit_layout, extract_circuit_location, extract_circuit_metrics
)

from ..utils.converters import datetime_to_ms, to_datetime

logger = logging.getLogger(__name__)

class FastF1Provider(Provider):
    _cache_initialized = False
    _memory_cache = {}

    def __init__(self):
        if not FastF1Provider._cache_initialized:
            self._setup_fastf1()
            FastF1Provider._cache_initialized = True

    def _get_from_cache(self, key: str):
        if key in self._memory_cache:
            value, expiry = self._memory_cache[key]
            if time.time() < expiry:
                logger.info(f"Memory Cache hit: {key}")
                return value
            else:
                logger.info(f"Memory Cache expired: {key}")
                del self._memory_cache[key]
        return None

    def _set_to_cache(self, key: str, value: Any, ttl: int = 600):
        self._memory_cache[key] = (value, time.time() + ttl)

    def _normalize_string(self, text: str) -> str:
        """Normalize string by removing accents and converting to lowercase."""
        if not text:
            return ""
        normalized = unicodedata.normalize('NFD', text)
        return "".join(c for c in normalized if unicodedata.category(c) != 'Mn').lower()

    def _setup_fastf1(self):
        # Configure FastF1 logging to reduce chatter
        fastf1_logger = logging.getLogger('fastf1')
        fastf1_logger.setLevel(logging.WARNING)
        
        # Configure cache
        cache_directory = os.path.join(os.getcwd(), '.cache')
        if not os.path.exists(cache_directory):
            os.makedirs(cache_directory)
        
        logger.info(f"Initializing FastF1 cache at {cache_directory}")
        fastf1.Cache.enable_cache(cache_directory)

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        cache_key = f"events_{year}"
        cached = self._get_from_cache(cache_key)
        if cached is not None:
            return cached

        logger.info(f"Entry: get_weekend_events(year={year})")
        try:
            schedule = fastf1.get_event_schedule(year)
            race_weekends = []
            for _, weekend in schedule.iterrows():
                race_weekend = extract_race_weekend(weekend)
                if race_weekend:
                    race_weekends.append(race_weekend)
            logger.info(f"Exit: get_weekend_events(year={year}) - Found {len(race_weekends)} race weekends")
            self._set_to_cache(cache_key, race_weekends, ttl=3600)
            return race_weekends
        except Exception:
            logger.exception(f"Error fetching weekend events for year {year}")
            return []

    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        cache_key = f"results_{year}_{round_number}_{session_type}"
        cached = self._get_from_cache(cache_key)
        if cached is not None:
            return cached

        logger.info(f"Entry: get_session_results(year={year}, round={round_number}, session={session_type})")
        try:
            session = fastf1.get_session(year, round_number, session_type)
            session.load(laps=True, telemetry=False,
                         weather=False, messages=True)

            if session.results is None or session.results.empty:
                logger.warning(f"No results found for {year} Round {round_number} ({session_type})")
                return None

            results = []
            for _, row in session.results.iterrows():
                driver_result = extract_driver_result(row, session, session_type)
                results.append(driver_result)

            logger.info(f"Exit: get_session_results(year={year}, round={round_number}, session={session_type}) - Found {len(results)} drivers")
            res = SessionResult(
                year=year,
                round=round_number,
                session_type=session_type,
                results=results
            )
            # Use short 1-minute TTL for memory cache to ensure live sessions update
            self._set_to_cache(cache_key, res, ttl=60)
            return res
        except Exception:
            logger.exception(f"Error fetching session results for {year} round {round_number} {session_type}")
            return None

    def get_drivers(self, year: int, round_number: int) -> List[DriverInfo]:
        cache_key = f"drivers_{year}_{round_number}"
        cached = self._get_from_cache(cache_key)
        if cached is not None:
            return cached

        logger.info(f"Entry: get_drivers(year={year}, round={round_number})")
        
        drivers = []
        # Look back through previous rounds of the current year
        for r in range(round_number - 1, 0, -1):
            logger.info(f"Attempting to fetch drivers from {year} Round {r} R")
            drivers = self._fetch_from_session(year, r, 'R')
            if drivers:
                break
        
        # If still no drivers, try the last round of the previous year
        if not drivers:
            logger.info(f"No drivers found in {year}, attempting last round of {year-1}")
            try:
                prev_schedule = fastf1.get_event_schedule(year - 1)
                if not prev_schedule.empty:
                    last_round = int(prev_schedule.iloc[-1]["RoundNumber"])
                    drivers = self._fetch_from_session(year - 1, last_round, 'R')
            except Exception:
                logger.exception(f"Error fetching drivers from previous year {year-1}")
            
        if drivers:
            logger.info(f"Exit: get_drivers(year={year}, round={round_number}) - Found {len(drivers)} drivers")
            self._set_to_cache(cache_key, drivers, ttl=3600)
            return drivers
            
        logger.warning(f"Exit: get_drivers(year={year}, round={round_number}) - No drivers found")
        return []

    def _fetch_from_session(self, year: int, round_number: int, session_type: str) -> List[DriverInfo]:
        """Internal helper to fetch drivers from a specific session."""
        logger.info(f"Entry: _fetch_from_session(year={year}, round={round_number}, session={session_type})")
        try:
            session = fastf1.get_session(year, round_number, session_type)
            session.load(laps=True, telemetry=False, weather=False, messages=False)
            
            if session.results is not None and not session.results.empty:
                drivers = []
                for _, row in session.results.iterrows():
                    driver_info = extract_driver_info(row)
                    if driver_info.id:
                        drivers.append(driver_info)
                logger.info(f"Exit: _fetch_from_session(year={year}, round={round_number}, session={session_type}) - Found {len(drivers)} drivers")
                return drivers
        except Exception as e:
            logger.exception(f"Could not load {session_type} session for drivers: {e}")
        
        logger.warning(f"Exit: _fetch_from_session(year={year}, round={round_number}, session={session_type}) - No drivers found")
        return []

    def get_circuit_data(self, year: int, round_number: int) -> Optional[Circuit]:
        cache_key = f"circuit_{year}_{round_number}"
        cached = self._get_from_cache(cache_key)
        if cached is not None:
            return cached

        logger.info(f"Entry: get_circuit_data(year={year}, round={round_number})")
        try:
            ergast = Ergast()
            ergast_circuits = ergast.get_circuits(season=year, round=round_number)
            location_info = extract_circuit_location(ergast_circuits)

            session = self._fetch_qualifying_session(year, round_number)
            metrics, layout, rotation = self._extract_circuit_metrics_and_layout(session)
            
            if (metrics["length_km"] == 0 or not layout) and location_info["location"] != "Unknown":
                metrics, layout, rotation = self._get_fallback_circuit_data(location_info, year, metrics, layout, rotation)

            event_date = to_datetime(session.event.EventDate)
            event_date_ms = datetime_to_ms(event_date) or 0

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
                rotation=rotation,
                max_speed_kmh=metrics["max_speed_kmh"],
                max_altitude_m=metrics["max_altitude_m"],
                min_altitude_m=metrics["min_altitude_m"]
            )
            logger.info(f"Exit: get_circuit_data(year={year}, round={round_number}) - Success for {circuit.circuit_name}")
            self._set_to_cache(cache_key, circuit, ttl=3600)
            return circuit

        except Exception:
            logger.exception(f"Error fetching circuit data for {year} Round {round_number}")
            return None

    def _fetch_qualifying_session(self, year: int, round_number: int) -> Any:
        logger.info(f"Entry: _fetch_qualifying_session(year={year}, round={round_number})")
        session = fastf1.get_session(year, round_number, 'Q')
        try:
            session.load(telemetry=True, laps=True, weather=False, messages=False)
        except Exception:
            logger.exception(f"Could not load telemetry for {year} Round {round_number}")
        
        logger.info(f"Exit: _fetch_qualifying_session(year={year}, round={round_number})")
        return session

    def _extract_circuit_metrics_and_layout(self, session: Any) -> tuple:
        metrics = extract_circuit_metrics(session)
        layout = extract_circuit_layout(session)
        rotation = 0.0
        try:
            circuit_info = session.get_circuit_info()
            if circuit_info is not None:
                rotation = circuit_info.rotation
        except Exception:
            logger.exception("Error extracting circuit info rotation")
        return metrics, layout, rotation

    def _get_fallback_circuit_data(self, location_info: dict, year: int, metrics: dict, layout: list, rotation: float) -> tuple:
        logger.info(f"Data missing for {year}, attempting historical fallback")
        historical_session = self._get_historical_session(
            location_info["location"], 
            location_info["country"], 
            year
        )
        if historical_session:
            logger.info(f"Applying historical fallback from {historical_session.event.EventDate.year}")
            return self._extract_circuit_metrics_and_layout(historical_session)
        return metrics, layout, rotation

    def _get_historical_session(self, location: str, country: str, current_year: int) -> Optional[Any]:
        logger.info(f"Entry: _get_historical_session(location={location}, country={country}, current_year={current_year})")
        
        norm_location = self._normalize_string(location)
        norm_country = self._normalize_string(country)
        
        for search_year in range(current_year - 1, current_year - 6, -1):
            try:
                schedule = fastf1.get_event_schedule(search_year)
                
                # Apply normalization to schedule for matching
                schedule_norm = schedule.copy()
                schedule_norm['LocationNorm'] = schedule_norm['Location'].apply(self._normalize_string)
                schedule_norm['CountryNorm'] = schedule_norm['Country'].apply(self._normalize_string)
                
                matches = schedule_norm[
                    (schedule_norm["LocationNorm"] == norm_location) & 
                    (schedule_norm["CountryNorm"] == norm_country)
                ]
                
                if not matches.empty:
                    historical_round_number = int(matches.iloc[-1]["RoundNumber"])
                    historical_session = fastf1.get_session(search_year, historical_round_number, "Q")
                    historical_session.load(telemetry=True, laps=True, weather=False, messages=False)
                    
                    if hasattr(historical_session, "laps") and not historical_session.laps.empty:
                        logger.info(f"Exit: _get_historical_session - Found valid historical data in {search_year}")
                        return historical_session
            except Exception:
                logger.exception(f"Error checking historical data for {search_year}")
                continue
        
        logger.info("Exit: _get_historical_session - No historical data found")
        return None
