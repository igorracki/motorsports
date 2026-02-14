import pandas as pd
import logging
from datetime import datetime
from typing import Any, Optional, List
from ..models import (
    DriverResult, DriverInfo, RaceDetails, QualifyingDetails,
    RaceWeekend, Session
)
from ..models.circuit import CircuitLayoutPoint
from .converters import get_scalar_value, to_milliseconds, to_datetime

logger = logging.getLogger(__name__)

# Session Type Constants
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


def extract_race_weekend(weekend: pd.Series) -> Optional[RaceWeekend]:
    sessions = []
    for session_index in range(1, 6):
        session_type = str(weekend.get(f"Session{session_index}", ""))
        local_time = to_datetime(weekend.get(f"Session{session_index}Date"))
        utc_time = to_datetime(weekend.get(f"Session{session_index}DateUtc"))

        if isinstance(local_time, datetime) and isinstance(utc_time, datetime):
            sessions.append(Session(
                type=session_type,
                time_local=local_time,
                time_utc=utc_time
            ))

    event_date = to_datetime(weekend.get('EventDate'))

    if isinstance(event_date, datetime):
        raw_round = weekend.get('RoundNumber')
        round_number = int(raw_round) if raw_round is not None and pd.notna(
            raw_round) else 0

        return RaceWeekend(
            round=round_number,
            full_name=str(weekend.get('OfficialEventName', "") or ""),
            name=str(weekend.get('EventName', "") or ""),
            location=str(weekend.get('Location', "") or ""),
            country=str(weekend.get('Country', "") or ""),
            start_date=event_date,
            sessions=sessions
        )
    return None


def extract_driver_result(row: pd.Series, session: Any, session_type: str) -> DriverResult:
    driver_code = str(get_scalar_value(row, 'Abbreviation') or "")
    driver_info = DriverInfo(
        number=driver_code,
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
            driver_laps_any: Any = session.laps.pick_drivers(driver_code)
            if hasattr(driver_laps_any, 'empty') and not driver_laps_any.empty:
                fastest_lap = driver_laps_any.pick_fastest()
                if fastest_lap is not None and not pd.isna(fastest_lap['LapTime']):
                    fastest_lap_ms = to_milliseconds(fastest_lap['LapTime'])
        except Exception as e:
            logger.warning(f"Could not extract fastest lap for driver {driver_code}: {e}")

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

    return driver_result


import math

def extract_circuit_location(ergast_circuits: pd.DataFrame) -> dict:
    """Extracts circuit location info from Ergast data."""
    location_info = {
        "circuit_name": "Unknown Circuit",
        "location": "Unknown",
        "country": "Unknown",
        "latitude": 0.0,
        "longitude": 0.0
    }

    if not ergast_circuits.empty:
        circuit_row = ergast_circuits.iloc[0]
        location_info["circuit_name"] = str(circuit_row.get('circuitName', "Unknown Circuit"))
        location_info["location"] = str(circuit_row.get('locality', "Unknown"))
        location_info["country"] = str(circuit_row.get('country', "Unknown"))
        location_info["latitude"] = float(circuit_row.get('lat', 0.0))
        location_info["longitude"] = float(circuit_row.get('long', 0.0))
    
    return location_info

def extract_circuit_metrics(session: Any) -> dict:
    """Extracts corners and length from session data."""
    metrics = {
        "corners": 0,
        "length_km": 0.0
    }

    try:
        circuit_info = session.get_circuit_info()
        if circuit_info is not None:
            metrics["corners"] = len(circuit_info.corners)
    except Exception as e:
        logger.warning(f"Could not get circuit info (corners): {e}")

    try:
        fastest_lap = session.laps.pick_fastest()
        telemetry = fastest_lap.get_telemetry()
        if not telemetry.empty:
            max_distance = telemetry['Distance'].max()
            metrics["length_km"] = round(float(max_distance) / 1000.0, 3)
    except Exception as e:
        logger.warning(f"Could not calculate circuit length: {e}")

    return metrics

def extract_circuit_layout(session: Any, rotation: float = 0.0) -> List[CircuitLayoutPoint]:
    """Extracts and rotates circuit layout points from telemetry."""
    try:
        if not hasattr(session, 'laps'):
            return []
            
        fastest_lap = session.laps.pick_fastest()
        telemetry = fastest_lap.get_telemetry()
        
        if telemetry.empty:
            return []

        # Downsample to approx 500 points
        total_points = len(telemetry)
        target_points = 500
        step = max(1, total_points // target_points)
        
        # Prepare rotation
        theta_rad = math.radians(rotation)
        cos_theta = math.cos(theta_rad)
        sin_theta = math.sin(theta_rad)

        layout_points = []
        for i in range(0, total_points, step):
            telemetry_row = telemetry.iloc[i]
            coord_x = float(telemetry_row['X'])
            coord_y = float(telemetry_row['Y'])
            
            # Apply rotation
            rotated_x = coord_x * cos_theta - coord_y * sin_theta
            rotated_y = coord_x * sin_theta + coord_y * cos_theta

            layout_points.append(CircuitLayoutPoint(
                x=rotated_x,
                y=-rotated_y # Invert Y for screen coordinates
            ))
            
        return layout_points
    except Exception as e:
        logger.error(f"Failed to extract circuit layout: {e}")
        return []
