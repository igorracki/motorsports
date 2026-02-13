from ..providers import Provider
from ..models import RaceWeekend, SessionResult
from typing import List, Optional


class F1Service:
    def __init__(self, provider: Provider):
        self.provider = provider

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        return self.provider.get_weekend_events(year)

    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        return self.provider.get_session_results(year, round_number, session_type)
