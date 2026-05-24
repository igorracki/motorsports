import logging
from ..providers import Provider
from ..models import RaceWeekend, SessionResult, DriverInfo
from ..models.circuit import Circuit
from typing import List, Optional

logger = logging.getLogger(__name__)


class F1Service:
    def __init__(self, provider: Provider):
        self.provider = provider

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        return self.provider.get_weekend_events(year)

    def get_session_results(self, year: int, round_number: int, session_type: str, force_reload: bool = False) -> Optional[SessionResult]:
        return self.provider.get_session_results(year, round_number, session_type, force_reload=force_reload)

    def get_circuit_data(self, year: int, round_number: int, force_reload: bool = False) -> Optional[Circuit]:
        return self.provider.get_circuit_data(year, round_number, force_reload=force_reload)

    def get_drivers(self, year: int, round_number: int, force_reload: bool = False) -> List[DriverInfo]:
        return self.provider.get_drivers(year, round_number, force_reload=force_reload)

    def clear_cache(self) -> bool:
        return self.provider.clear_cache()
