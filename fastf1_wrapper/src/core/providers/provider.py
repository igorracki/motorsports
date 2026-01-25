from abc import ABC, abstractmethod
from ..models import RaceResult, Event
from typing import List


class Provider(ABC):
    @abstractmethod
    def get_race_result(self, year: int, round: int, session_type: str) -> RaceResult:
        pass

    @abstractmethod
    def get_weekend_events(self, year: int) -> List[Event]:
        pass
