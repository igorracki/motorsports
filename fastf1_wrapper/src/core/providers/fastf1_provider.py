import fastf1
from fastf1.ergast import Ergast
import pandas as pd
import os
import logging
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

    def get_drivers(self, year: int, round_number: int) -> List[DriverInfo]:
        logger.info(f"Entry: get_drivers(year={year}, round={round_number})")
        try:
            # Primary attempt: Use 'R' (Race) session to get the entry list/results
            session = fastf1.get_session(year, round_number, 'R')
            
            try:
                # Use light load to get results/entry list
                session.load(laps=False, telemetry=False, weather=False, messages=False)
            except Exception as e:
                logger.warning(f"Could not load session for drivers, trying fallback: {e}")
            
            # If session results are available, use them
            if session.results is not None and not session.results.empty:
                drivers = []
                for _, row in session.results.iterrows():
                    driver_info = extract_driver_info(row)
                    if driver_info.id:
                        drivers.append(driver_info)
                
                if drivers:
                    logger.info(f"Exit: get_drivers(year={year}, round={round_number}) - Found {len(drivers)} drivers via FastF1")
                    return drivers
            
            # Fallback: Use Ergast for season driver info if session results are missing (e.g. future seasons)
            logger.info(f"Falling back to Ergast season driver info for year {year}")
            ergast = Ergast()
            driver_info_response = ergast.get_driver_info(season=year)
            
            # Ergast responses in FastF1 3.x+ return an ErgastRawResponse which has a 'content' attribute
            # containing the list of dataframes.
            if not driver_info_response.content:
                logger.warning(f"No content found in Ergast response for year {year}")
                return []
                
            driver_info_df = driver_info_response.content[0]
            
            if driver_info_df.empty:
                logger.warning(f"No drivers found via Ergast for year {year}")
                return []
                
            drivers = []
            for _, row in driver_info_df.iterrows():
                # Prefer abbreviation if available, otherwise use driverId (slug)
                abbr = row.get('abbreviation')
                driver_id = str(abbr if abbr else row.get('driverId')).upper()
                
                # Ensure driver number is an integer string (avoiding 1.0)
                raw_number = row.get('driverNumber')
                try:
                    if raw_number and not pd.isna(raw_number):
                        driver_number = str(int(float(raw_number)))
                    else:
                        driver_number = "0"
                except (ValueError, TypeError):
                    driver_number = str(raw_number or "0")
                
                drivers.append(DriverInfo(
                    id=driver_id,
                    number=driver_number,
                    full_name=f"{row.get('givenName')} {row.get('familyName')}",
                    country_code=str(row.get('nationality') or ""),
                    team_name="" # Ergast season driver info doesn't reliably map to constructors
                ))
            
            logger.info(f"Exit: get_drivers(year={year}, round={round_number}) - Found {len(drivers)} drivers via Ergast")
            return drivers
        except Exception:
            logger.exception(f"Error fetching drivers for {year} round {round_number}")
            return []

    def get_circuit_data(self, year: int, round_number: int) -> Optional[Circuit]:
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
            return circuit

        except Exception:
            logger.exception(f"Error fetching circuit data for {year} Round {round_number}")
            return None

    def _fetch_qualifying_session(self, year: int, round_number: int) -> Any:
        session = fastf1.get_session(year, round_number, 'Q')
        try:
            session.load(telemetry=True, laps=True, weather=False, messages=False)
        except Exception:
            logger.warning(f"Could not load telemetry for {year} Round {round_number}")
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
            pass
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
        for search_year in range(current_year - 1, current_year - 6, -1):
            try:
                schedule = fastf1.get_event_schedule(search_year)
                matches = schedule[(schedule["Location"] == location) & (schedule["Country"] == country)]
                if not matches.empty:
                    historical_round_number = int(matches.iloc[-1]["RoundNumber"])
                    historical_session = fastf1.get_session(search_year, historical_round_number, "Q")
                    historical_session.load(telemetry=True, laps=True, weather=False, messages=False)
                    
                    if hasattr(historical_session, "laps") and not historical_session.laps.empty:
                        logger.info(f"Exit: _get_historical_session - Found valid historical data in {search_year}")
                        return historical_session
            except Exception:
                continue
        
        logger.info(f"Exit: _get_historical_session - No historical data found")
        return None
