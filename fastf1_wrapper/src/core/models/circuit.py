from dataclasses import dataclass
from typing import List, Optional

@dataclass
class CircuitLayoutPoint:
    x: float
    y: float

@dataclass
class Circuit:
    circuit_name: str
    location: str
    country: str
    latitude: float
    longitude: float
    length_km: float
    corners: int
    layout: List[CircuitLayoutPoint]
    event_name: str
    event_date_ms: int
    rotation: float = 0.0
    max_speed_kmh: Optional[float] = None
    max_altitude_m: Optional[float] = None
    min_altitude_m: Optional[float] = None
