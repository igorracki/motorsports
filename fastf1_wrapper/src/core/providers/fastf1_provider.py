import fastf1
from .provider import Provider
from ..models import Session, Event
from typing import List


class FastF1Provider(Provider):

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
