from ..providers import Provider


class SessionService:

    def __init__(self, provider: Provider):
        self.provider = provider

    def get_race_result(self, year: int, round: int, session_type: str) -> RaceResult:
        return self.provider.get_race_result(year, round, session_type)
