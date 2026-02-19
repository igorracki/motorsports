from dataclasses import dataclass, field
from typing import List, Optional


@dataclass
class DriverInfo:
    id: str
    number: str
    full_name: str
    country_code: str
    team_name: str


@dataclass
class RaceDetails:
    grid_position: int
    status: str
    positions_change: int = 0


@dataclass
class QualifyingDetails:
    q1_ms: Optional[int] = None
    q2_ms: Optional[int] = None
    q3_ms: Optional[int] = None


@dataclass
class DriverResult:
    position: int
    driver: DriverInfo
    laps: int
    status: str
    total_time_ms: Optional[int] = None
    gap_ms: Optional[int] = None
    fastest_lap_ms: Optional[int] = None
    race_details: Optional[RaceDetails] = None
    qualifying_details: Optional[QualifyingDetails] = None


@dataclass
class SessionResult:
    year: int
    round: int
    session_type: str
    results: List[DriverResult] = field(default_factory=list)
