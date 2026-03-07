import pandas as pd
import logging
from typing import Any
from ..models import DriverResult, DriverInfo, RaceDetails, QualifyingDetails
from .converters import get_scalar_value, to_milliseconds

logger = logging.getLogger(__name__)

SESSION_RACE = 'R'
SESSION_RACE_FULL = 'RACE'
SESSION_SPRINT = 'S'
SESSION_SPRINT_FULL = 'SPRINT'
SESSION_QUALIFYING = 'Q'
SESSION_QUALIFYING_FULL = 'QUALIFYING'
SESSION_SPRINT_QUALIFYING = 'SQ'
SESSION_SPRINT_QUALIFYING_FULL = 'SPRINT QUALIFYING'

RACE_SESSION_TYPES = {SESSION_RACE, SESSION_RACE_FULL, SESSION_SPRINT, SESSION_SPRINT_FULL}
QUALIFYING_SESSION_TYPES = {
    SESSION_QUALIFYING, SESSION_QUALIFYING_FULL,
    SESSION_SPRINT_QUALIFYING, SESSION_SPRINT_QUALIFYING_FULL
}

def extract_driver_info(row: pd.Series) -> DriverInfo:
    id = str(get_scalar_value(row, 'Abbreviation') or "")
    
    # Ensure driver number is an integer string (avoiding 1.0)
    raw_number = get_scalar_value(row, 'DriverNumber')
    try:
        if raw_number and not pd.isna(raw_number):
            driver_number = str(int(float(raw_number)))
        else:
            driver_number = "0"
    except (ValueError, TypeError):
        driver_number = str(raw_number or "0")

    return DriverInfo(
        id=id,
        number=driver_number,
        full_name=str(get_scalar_value(row, 'FullName') or ""),
        country_code=str(get_scalar_value(row, 'CountryCode') or ""),
        team_name=str(get_scalar_value(row, 'TeamName') or "")
    )

def extract_driver_result(row: pd.Series, session: Any, session_type: str) -> DriverResult:
    driver_info = extract_driver_info(row)
    id = driver_info.id
    logger.info(f"Extracting result for driver: {id}")

    status = str(get_scalar_value(row, 'Status') or "")

    laps_val = get_scalar_value(row, 'Laps')
    if laps_val is None or pd.isna(laps_val):
        laps_val = get_scalar_value(row, 'NoLaps')

    laps_completed = int(
        laps_val) if laps_val is not None and pd.notna(laps_val) else 0
    
    # Position can be float or int in FastF1, handle carefully
    raw_position = get_scalar_value(row, 'Position')
    if pd.notna(raw_position):
        try:
            finish_position = int(float(raw_position))
        except (ValueError, TypeError):
            finish_position = 0
    else:
        finish_position = 0

    raw_time_delta = get_scalar_value(row, 'Time')
    total_time_ms = None
    gap_ms = None

    normalized_session_type = session_type.upper()

    if normalized_session_type in RACE_SESSION_TYPES:
        if finish_position == 1:
            total_time_ms = to_milliseconds(raw_time_delta)
            gap_ms = 0
        elif raw_time_delta is not None and not pd.isna(raw_time_delta):
            gap_ms = to_milliseconds(raw_time_delta)
    else:
        total_time_ms = to_milliseconds(raw_time_delta)

    fastest_lap_ms = None
    if hasattr(session, 'laps'):
        try:
            driver_laps_any: Any = session.laps.pick_drivers(id)
            if hasattr(driver_laps_any, 'empty') and not driver_laps_any.empty:
                fastest_lap = driver_laps_any.pick_fastest()
                if fastest_lap is not None and not pd.isna(fastest_lap['LapTime']):
                    fastest_lap_ms = to_milliseconds(fastest_lap['LapTime'])
        except Exception:
            logger.exception(f"Could not extract fastest lap for driver {id}")

    driver_result = DriverResult(
        position=finish_position,
        driver=driver_info,
        laps=laps_completed,
        status=status,
        total_time_ms=total_time_ms,
        gap_ms=gap_ms,
        fastest_lap_ms=fastest_lap_ms
    )

    if normalized_session_type in RACE_SESSION_TYPES:
        grid_position = get_scalar_value(row, 'GridPosition')
        driver_result.race_details = RaceDetails(
            grid_position=int(grid_position) if pd.notna(
                grid_position) else 0,
            status=status,
            positions_change=0
        )

    elif normalized_session_type in QUALIFYING_SESSION_TYPES:
        # For Qualifying and Sprint Qualifying, extract segment times.
        # FastF1 uses Q1/Q2/Q3 for Qualifying and sometimes SQ1/SQ2/SQ3 for Sprint Qualifying.
        # We try both sets of column names.
        q1_val = get_scalar_value(row, 'Q1')
        if pd.isna(q1_val): q1_val = get_scalar_value(row, 'SQ1')
        q1_ms = to_milliseconds(q1_val)

        q2_val = get_scalar_value(row, 'Q2')
        if pd.isna(q2_val): q2_val = get_scalar_value(row, 'SQ2')
        q2_ms = to_milliseconds(q2_val)

        q3_val = get_scalar_value(row, 'Q3')
        if pd.isna(q3_val): q3_val = get_scalar_value(row, 'SQ3')
        q3_ms = to_milliseconds(q3_val)

        driver_result.qualifying_details = QualifyingDetails(
            q1_ms=q1_ms,
            q2_ms=q2_ms,
            q3_ms=q3_ms
        )
        
        # Ensure we have a representative best lap for the summary table.
        # Primary: The absolute fastest lap recorded in the session.
        # Secondary: The time from the latest qualifying segment reached.
        if driver_result.fastest_lap_ms is None:
            if q3_ms is not None: driver_result.fastest_lap_ms = q3_ms
            elif q2_ms is not None: driver_result.fastest_lap_ms = q2_ms
            elif q1_ms is not None: driver_result.fastest_lap_ms = q1_ms

    logger.info(f"Successfully extracted result for driver: {id} (Pos: {finish_position})")
    return driver_result
