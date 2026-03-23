import pandas as pd
import logging
from typing import Any, List
from ..models.circuit import CircuitLayoutPoint

logger = logging.getLogger(__name__)

def extract_circuit_location(ergast_circuits: Any) -> dict:
    logger.info("Entry: extract_circuit_location")
    location_info = {
        "circuit_name": "Unknown Circuit",
        "location": "Unknown",
        "country": "Unknown",
        "latitude": 0.0,
        "longitude": 0.0
    }

    if ergast_circuits is not None and hasattr(ergast_circuits, 'empty') and not ergast_circuits.empty:
        circuit_row = ergast_circuits.iloc[0]
        location_info["circuit_name"] = str(circuit_row.get('circuitName', "Unknown Circuit"))
        location_info["location"] = str(circuit_row.get('locality', "Unknown"))
        location_info["country"] = str(circuit_row.get('country', "Unknown"))
        location_info["latitude"] = float(circuit_row.get('lat', 0.0))
        location_info["longitude"] = float(circuit_row.get('long', 0.0))
    
    logger.info(f"Exit: extract_circuit_location - Success for {location_info['circuit_name']}")
    return location_info

def extract_circuit_metrics(session: Any) -> dict:
    logger.info(f"Entry: extract_circuit_metrics(session={session.event.EventName})")
    metrics = {
        "corners": 0,
        "length_km": 0.0,
        "max_speed_kmh": 0.0,
        "max_altitude_m": 0.0,
        "min_altitude_m": 0.0
    }

    try:
        # Check if session info is loaded before calling get_circuit_info
        if hasattr(session, '_session_info') and session._session_info is not None:
            circuit_info = session.get_circuit_info()
            if circuit_info is not None:
                metrics["corners"] = len(circuit_info.corners)
    except Exception:
        logger.warning("Could not get circuit info (corners) - data may not be loaded")

    try:
        # Check if laps are loaded
        if hasattr(session, '_laps') and session._laps is not None and not session._laps.empty:
            fastest_lap = session.laps.pick_fastest()
            if fastest_lap is not None:
                telemetry = fastest_lap.get_telemetry()
                if not telemetry.empty:
                    max_distance = telemetry['Distance'].max()
                    metrics["length_km"] = float(max_distance) / 1000.0
                    
                    # Additional metrics
                    metrics["max_speed_kmh"] = float(telemetry['Speed'].max())
                    metrics["max_altitude_m"] = float(telemetry['Z'].max())
                    metrics["min_altitude_m"] = float(telemetry['Z'].min())
    except Exception:
        logger.warning("Could not calculate circuit metrics from telemetry - data may not be loaded")

    logger.info(f"Exit: extract_circuit_metrics - {metrics['corners']} corners, {metrics['length_km']:.3f} km, {metrics['max_speed_kmh']} km/h")
    return metrics

def extract_circuit_layout(session: Any) -> List[CircuitLayoutPoint]:
    logger.info(f"Entry: extract_circuit_layout(session={session.event.EventName})")
    try:
        # Check if laps are loaded
        if not hasattr(session, '_laps') or session._laps is None or session._laps.empty:
            return []
            
        fastest_lap = session.laps.pick_fastest()
        if fastest_lap is None:
            return []
            
        telemetry = fastest_lap.get_telemetry()
        
        if telemetry.empty:
            return []

        total_points = len(telemetry)
        target_points = 500
        step = max(1, total_points // target_points)
        
        layout_points = []
        for i in range(0, total_points, step):
            telemetry_row = telemetry.iloc[i]
            x = float(telemetry_row['X'])
            y = float(telemetry_row['Y'])
            
            layout_points.append(CircuitLayoutPoint(
                x=x,
                y=y
            ))
            
        logger.info(f"Exit: extract_circuit_layout - Extracted {len(layout_points)} points")
        return layout_points
    except Exception:
        logger.exception("Failed to extract circuit layout")
        return []
