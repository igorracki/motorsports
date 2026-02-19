import pandas as pd
import logging
from typing import Any
from ..models import DriverResult, DriverInfo, RaceDetails, QualifyingDetails
from .converters import get_scalar_value, to_milliseconds

logger = logging.getLogger(__name__)

SESSION_RACE = 'R'
SESSION_RACE_FULL = 'RACE'
SESSION_QUALIFYING = 'Q'
SESSION_QUALIFYING_FULL = 'QUALIFYING'
SESSION_SPRINT_QUALIFYING = 'SQ'
SESSION_SPRINT_QUALIFYING_FULL = 'SPRINT QUALIFYING'

RACE_SESSION_TYPES = {SESSION_RACE, SESSION_RACE_FULL}
QUALIFYING_SESSION_TYPES = {
    SESSION_QUALIFYING, SESSION_QUALIFYING_FULL,
    SESSION_SPRINT_QUALIFYING, SESSION_SPRINT_QUALIFYING_FULL
}

def extract_driver_result(row: pd.Series, session: Any, session_type: str) -> DriverResult:
    id = str(get_scalar_value(row, 'Abbreviation') or "")
    logger.info(f"Extracting result for driver: {id}")
    driver_info = DriverInfo(
        id=id,
        number=str(get_scalar_value(row, 'DriverNumber') or ""),
        full_name=str(get_scalar_value(row, 'FullName') or ""),
        country_code=str(get_scalar_value(row, 'CountryCode') or ""),
        team_name=str(get_scalar_value(row, 'TeamName') or "")
    )

    status = str(get_scalar_value(row, 'Status') or "")

    laps_val = get_scalar_value(row, 'Laps')
    if laps_val is None or pd.isna(laps_val):
        laps_val = get_scalar_value(row, 'NoLaps')

    laps_completed = int(
        laps_val) if laps_val is not None and pd.notna(laps_val) else 0
    finish_position = int(get_scalar_value(row, 'Position')) if pd.notna(
        get_scalar_value(row, 'Position')) else 0

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
        driver_result.qualifying_details = QualifyingDetails(
            q1_ms=to_milliseconds(get_scalar_value(row, 'Q1')),
            q2_ms=to_milliseconds(get_scalar_value(row, 'Q2')),
            q3_ms=to_milliseconds(get_scalar_value(row, 'Q3'))
        )

    logger.info(f"Successfully extracted result for driver: {id} (Pos: {finish_position})")
    return driver_result
