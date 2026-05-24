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
        fastf1_logger = logging.getLogger('fastf1')
        fastf1_logger.setLevel(logging.WARNING)
        
        cache_directory = os.environ.get('FASTF1_CACHE_DIR', os.path.join(os.getcwd(), '.cache'))
        
        try:
            if not os.path.exists(cache_directory):
                logger.info(f"Creating cache directory at {cache_directory}")
                os.makedirs(cache_directory, exist_ok=True)
            
            test_file = os.path.join(cache_directory, '.write_test')
            with open(test_file, 'w') as f:
                f.write('test')
            os.remove(test_file)
            
            logger.info(f"Initializing FastF1 cache at {cache_directory}")
            fastf1.Cache.enable_cache(cache_directory)
        except Exception as e:
            logger.error(f"Failed to initialize persistent cache at {cache_directory}: {e}")
            logger.warning("Falling back to memory-only cache or temporary directory")
            # FastF1 still works without persistent cache, but it's slower.
            # We don't call enable_cache if it fails.

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        cache_key = f"events_{year}"
        cached = self._get_from_cache(cache_key)
        if cached is not None:
            return cached

        try:
            schedule = fastf1.get_event_schedule(year)
            race_weekends = []
            for _, weekend in schedule.iterrows():
                race_weekend = extract_race_weekend(weekend)
                if race_weekend:
                    race_weekends.append(race_weekend)
            self._set_to_cache(cache_key, race_weekends, ttl=3600)
            return race_weekends
        except Exception:
            logger.exception(f"Error fetching weekend events for year {year}")
            return []

    def get_session_results(self, year: int, round_number: int, session_type: str, force_reload: bool = False) -> Optional[SessionResult]:
        cache_key = f"results_{year}_{round_number}_{session_type}"
        if not force_reload:
            cached = self._get_from_cache(cache_key)
            if cached is not None:
                return cached

        try:
            session = fastf1.get_session(year, round_number, session_type)
            session.load(laps=True, telemetry=False,
                         weather=False, messages=True,
                         force_reload=force_reload)

            if session.results is None or session.results.empty:
                logger.warning(f"No results found for {year} Round {round_number} ({session_type})")
                return None

            results = []
            for _, row in session.results.iterrows():
                driver_result = extract_driver_result(row, session, session_type)
                results.append(driver_result)

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

    def get_drivers(self, year: int, round_number: int, force_reload: bool = False) -> List[DriverInfo]:
        cache_key = f"drivers_{year}_{round_number}"
        if not force_reload:
            cached = self._get_from_cache(cache_key)
            if cached is not None:
                return cached
        
        drivers = []
        
        # 1. Try current round FP1 first - it often has the entry list even if no results
        drivers = self._fetch_from_session(year, round_number, 'FP1', force_reload=force_reload)
        
        # 2. Look back through previous rounds of the current year
        if not drivers:
            for r in range(round_number - 1, 0, -1):
                drivers = self._fetch_from_session(year, r, 'R', force_reload=force_reload)
                if drivers:
                    break
        
        # 3. If still no drivers, try the last round of the previous year
        if not drivers:
            try:
                prev_schedule = fastf1.get_event_schedule(year - 1)
                if not prev_schedule.empty:
                    last_round = int(prev_schedule.iloc[-1]["RoundNumber"])
                    drivers = self._fetch_from_session(year - 1, last_round, 'R', force_reload=force_reload)
            except Exception:
                logger.exception(f"Error fetching drivers from previous year {year-1}")
            
        if drivers:
            self._set_to_cache(cache_key, drivers, ttl=3600)
            return drivers
            
        return []

    def _fetch_from_session(self, year: int, round_number: int, session_type: str, force_reload: bool = False) -> List[DriverInfo]:
        """Internal helper to fetch drivers from a specific session."""
        try:
            session = fastf1.get_session(year, round_number, session_type)
            # Load without telemetry/laps to be faster and more resilient
            session.load(laps=False, telemetry=False, weather=False, messages=False,
                         force_reload=force_reload)
            
            drivers = []
            # Try to get from results first (Ergast)
            if session.results is not None and not session.results.empty:
                for _, row in session.results.iterrows():
                    driver_info = extract_driver_info(row)
                    if driver_info.id:
                        drivers.append(driver_info)
            
            # If no results, try the entry list (F1 API via FastF1)
            # Use internal attribute check to avoid triggering DataNotLoadedError property getter
            if not drivers and getattr(session, '_drivers', None) is not None:
                for driver_number in session.drivers:
                    try:
                        driver_data = session.get_driver(driver_number)
                        driver_info = extract_driver_info(driver_data)
                        if driver_info.id:
                            drivers.append(driver_info)
                    except Exception:
                        continue

            if drivers:
                return drivers
        except Exception as e:
            logger.warning(f"Could not load {session_type} session for drivers: {e}")
        
        return []

    def get_circuit_data(self, year: int, round_number: int, force_reload: bool = False) -> Optional[Circuit]:
        cache_key = f"circuit_{year}_{round_number}"
        if not force_reload:
            cached = self._get_from_cache(cache_key)
            if cached is not None:
                return cached

        try:
            ergast = Ergast()
            ergast_circuits = ergast.get_circuits(season=year, round=round_number)
            location_info = extract_circuit_location(ergast_circuits)

            session = self._fetch_qualifying_session(year, round_number, force_reload=force_reload)
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
            self._set_to_cache(cache_key, circuit, ttl=3600)
            return circuit

        except Exception:
            logger.exception(f"Error fetching circuit data for {year} Round {round_number}")
            return None

    def _fetch_qualifying_session(self, year: int, round_number: int, force_reload: bool = False) -> Any:
        session = fastf1.get_session(year, round_number, 'Q')
        try:
            session.load(telemetry=True, laps=True, weather=False, messages=False,
                         force_reload=force_reload)
        except Exception:
            logger.exception(f"Could not load telemetry for {year} Round {round_number}")
        
        return session

    def _extract_circuit_metrics_and_layout(self, session: Any) -> tuple:
        metrics = extract_circuit_metrics(session)
        layout = extract_circuit_layout(session)
        rotation = 0.0
        try:
            # Check if session info is loaded before calling get_circuit_info
            if getattr(session, '_session_info', None) is not None:
                circuit_info = session.get_circuit_info()
                if circuit_info is not None:
                    rotation = circuit_info.rotation
        except Exception:
            logger.warning("Error extracting circuit info rotation - data may not be loaded")
        return metrics, layout, rotation

    def _get_fallback_circuit_data(self, location_info: dict, year: int, metrics: dict, layout: list, rotation: float) -> tuple:
        historical_session = self._get_historical_session(
            location_info["location"], 
            location_info["country"], 
            year
        )
        if historical_session:
            return self._extract_circuit_metrics_and_layout(historical_session)
        return metrics, layout, rotation

    def _get_historical_session(self, location: str, country: str, current_year: int) -> Optional[Any]:
        norm_location = self._normalize_string(location)
        norm_country = self._normalize_string(country)
        
        # Address Ergast vs FastF1 string mismatches
        country_aliases = {
            'usa': 'united states',
            'uk': 'united kingdom',
            'uae': 'united arab emirates'
        }
        location_aliases = {
            'monte carlo': 'monaco',
            'yas marina': 'yas island',
            'spa': 'spa-francorchamps',
        }
        
        search_country = country_aliases.get(norm_country, norm_country)
        search_location = location_aliases.get(norm_location, norm_location)
        
        for search_year in range(current_year - 1, current_year - 6, -1):
            try:
                schedule = fastf1.get_event_schedule(search_year)
                
                # Apply normalization to schedule for matching
                schedule_norm = schedule.copy()
                schedule_norm['LocationNorm'] = schedule_norm['Location'].apply(self._normalize_string)
                schedule_norm['CountryNorm'] = schedule_norm['Country'].apply(self._normalize_string)
                schedule_norm['EventNameNorm'] = schedule_norm['EventName'].apply(self._normalize_string)
                
                # 1. Filter by Country
                country_matches = schedule_norm[schedule_norm['CountryNorm'] == search_country]
                
                if country_matches.empty:
                    continue
                    
                # 2. Filter by Location (exact, contains, or in Event Name)
                loc_matches = country_matches[
                    (country_matches['LocationNorm'] == search_location) |
                    (country_matches['LocationNorm'].str.contains(search_location, na=False)) |
                    (country_matches['EventNameNorm'].str.contains(search_location, na=False))
                ]
                
                if not loc_matches.empty:
                    historical_round_number = int(loc_matches.iloc[-1]["RoundNumber"])
                    
                    # Try Qualifying first, then Race as fallback for circuit data
                    for session_type in ["Q", "R"]:
                        historical_session = fastf1.get_session(search_year, historical_round_number, session_type)
                        
                        try:
                            historical_session.load(telemetry=True, laps=True, weather=False, messages=False)
                        except Exception as e:
                            logger.warning(f"Failed to load historical session {search_year} {session_type}: {e}")
                            continue
                        
                        if getattr(historical_session, "_laps", None) is not None and not historical_session._laps.empty:
                            return historical_session
            except Exception:
                logger.warning(f"Error checking historical data for {search_year}", exc_info=True)
                continue
        
        return None

    def clear_cache(self) -> bool:
        pass
