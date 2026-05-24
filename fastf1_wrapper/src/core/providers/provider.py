from abc import ABC, abstractmethod
from ..models import RaceWeekend, SessionResult, DriverInfo
from ..models.circuit import Circuit
from typing import List, Optional


class Provider(ABC):
    @abstractmethod
    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        pass

    @abstractmethod
    def get_session_results(self, year: int, round_number: int, session_type: str, force_reload: bool = False) -> Optional[SessionResult]:
        pass

    @abstractmethod
    def get_circuit_data(self, year: int, round_number: int, force_reload: bool = False) -> Optional[Circuit]:
        pass

    @abstractmethod
    def get_drivers(self, year: int, round_number: int, force_reload: bool = False) -> List[DriverInfo]:
        pass

    @abstractmethod
    def clear_cache(self) -> bool:
        pass
