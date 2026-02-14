from abc import ABC, abstractmethod
from ..models import RaceWeekend, SessionResult
from ..models.circuit import Circuit
from typing import List, Optional


class Provider(ABC):
    @abstractmethod
    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        pass

    @abstractmethod
    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        pass

    @abstractmethod
    def get_circuit_data(self, year: int, round_number: int) -> Optional[Circuit]:
        pass
