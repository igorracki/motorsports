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

    if not ergast_circuits.empty:
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
        "length_km": 0.0
    }

    try:
        circuit_info = session.get_circuit_info()
        if circuit_info is not None:
            metrics["corners"] = len(circuit_info.corners)
    except Exception:
        logger.exception("Could not get circuit info (corners)")

    try:
        fastest_lap = session.laps.pick_fastest()
        telemetry = fastest_lap.get_telemetry()
        if not telemetry.empty:
            max_distance = telemetry['Distance'].max()
            metrics["length_km"] = float(max_distance) / 1000.0
    except Exception:
        logger.exception("Could not calculate circuit length")

    logger.info(f"Exit: extract_circuit_metrics - {metrics['corners']} corners, {metrics['length_km']:.3f} km")
    return metrics

def extract_circuit_layout(session: Any) -> List[CircuitLayoutPoint]:
    logger.info(f"Entry: extract_circuit_layout(session={session.event.EventName})")
    try:
        if not hasattr(session, 'laps'):
            return []
            
        fastest_lap = session.laps.pick_fastest()
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
