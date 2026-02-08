from abc import ABC, abstractmethod
from ..models import Event
from typing import List


class Provider(ABC):
    @abstractmethod
    def get_weekend_events(self, year: int) -> List[Event]:
        pass
