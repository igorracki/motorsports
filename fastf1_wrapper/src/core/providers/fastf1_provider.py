import fastf1
from .provider import Provider
from ..models import RaceResult, SessionInfo, DriverResult, Driver, Session, Event
from ..utils import format_lap_time
from typing import List


class FastF1Provider(Provider):
    def get_race_result(self, year: int, round: int, session_type: str) -> RaceResult:
        session = fastf1.get_session(year, round, session_type)
        session.load()
        session_info = SessionInfo(
            name=session.event['EventName'],
            session_type=session.name,
            date=session.event['EventDate'],
            event_name=session.event['OfficialEventName'],
            location=session.event['Location'],
            country=session.event['Country']
        )

        laps = session.laps
        driver_results = []
        for _, row in session.results.iterrows():
            driver_number = int(row['DriverNumber'])
            driver_laps = laps.pick_drivers(driver_number)
            driver_fastest_lap = driver_laps.pick_fastest()
            fastest_lap_time = driver_fastest_lap['LapTime']
            position_delta = int(row['GridPosition']) - int(row['Position'])

            driver = Driver(
                driver_number=driver_number,
                first_name=row['FirstName'],
                surname=row['LastName'],
                team=row['TeamName']
            )

            driver_result = DriverResult(
                position=int(row['Position']),
                driver=driver,
                fastest_lap_time=format_lap_time(fastest_lap_time),
                positions_gained=position_delta
            )

            driver_results.append(driver_result)

        return RaceResult(
            session_info=session_info,
            results=driver_results
        )

    def get_weekend_events(self, year: int) -> List[Event]:
        schedule = fastf1.get_event_schedule(year)
        events = []
        for _, weekend in schedule.iterrows():
            sessions = []
            for n in range(1, 6):
                session = Session(
                    type=weekend[f"Session{n}"],
                    time_local=weekend[f"Session{n}Date"],
                    time_utc=weekend[f"Session{n}DateUtc"]
                )
                sessions.append(session)

            event = Event(
                round=weekend['RoundNumber'],
                full_name=weekend['OfficialEventName'],
                name=weekend['EventName'],
                location=weekend['Location'],
                country=weekend['Country'],
                start_date=weekend['EventDate'],
                sessions=sessions
            )
            events.append(event)

        return events
