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
    event_date: str
    # lap_record: Optional[LapRecord] = None
