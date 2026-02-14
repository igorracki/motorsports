import logging
from ..providers import Provider
from ..models import RaceWeekend, SessionResult
from ..models.circuit import Circuit
from typing import List, Optional

logger = logging.getLogger(__name__)


class F1Service:
    def __init__(self, provider: Provider):
        self.provider = provider

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        logger.info(f"Entry: get_weekend_events(year={year})")
        results = self.provider.get_weekend_events(year)
        logger.info(f"Exit: get_weekend_events(year={year}) - Found {len(results)} events")
        return results

    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        logger.info(f"Entry: get_session_results(year={year}, round={round_number}, session={session_type})")
        result = self.provider.get_session_results(year, round_number, session_type)
        count = len(result.results) if result and result.results else 0
        logger.info(f"Exit: get_session_results(year={year}, round={round_number}, session={session_type}) - Found {count} results")
        return result

    def get_circuit_data(self, year: int, round_number: int) -> Optional[Circuit]:
        logger.info(f"Entry: get_circuit_data(year={year}, round={round_number})")
        result = self.provider.get_circuit_data(year, round_number)
        logger.info(f"Exit: get_circuit_data(year={year}, round={round_number}) - Success: {result is not None}")
        return result
