from dataclasses import dataclass


@dataclass
class Session:
    type: str
    time_utc_ms: int
    utc_offset_ms: int
