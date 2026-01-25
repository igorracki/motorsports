from ..providers import Provider
from ..models import Event
from typing import List


class EventService:
    def __init__(self, provider: Provider):
        self.provider = provider

    def get_weekend_events(self, year: int) -> List[Event]:
        return self.provider.get_weekend_events(year)
