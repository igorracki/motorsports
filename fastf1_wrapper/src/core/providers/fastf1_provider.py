import fastf1
import os
import pandas as pd
from datetime import datetime
from typing import List, Optional, Any
from .provider import Provider
from ..models import (
    Session, RaceWeekend, SessionResult, DriverResult,
    DriverInfo, RaceDetails, QualifyingDetails
)
from ..utils.converters import to_milliseconds, get_scalar_value

cache_directory = os.path.join(os.getcwd(), '.cache')
if not os.path.exists(cache_directory):
    os.makedirs(cache_directory)
fastf1.Cache.enable_cache(cache_directory)


class FastF1Provider(Provider):

    def get_weekend_events(self, year: int) -> List[RaceWeekend]:
        schedule = fastf1.get_event_schedule(year)
        race_weekends = []
        for _, weekend in schedule.iterrows():
            sessions = []
            for session_index in range(1, 6):
                session_type = str(weekend.get(f"Session{session_index}", ""))
                local_time = weekend.get(f"Session{session_index}Date")
                utc_time = weekend.get(f"Session{session_index}DateUtc")

                # Convert pandas Timestamps to standard python datetime
                if local_time is not None and hasattr(local_time, 'to_pydatetime'):
                    local_time = local_time.to_pydatetime()
                if utc_time is not None and hasattr(utc_time, 'to_pydatetime'):
                    utc_time = utc_time.to_pydatetime()

                if isinstance(local_time, datetime) and isinstance(utc_time, datetime):
                    sessions.append(Session(
                        type=session_type,
                        time_local=local_time,
                        time_utc=utc_time
                    ))

            event_date = weekend.get('EventDate')
            if event_date is not None and hasattr(event_date, 'to_pydatetime'):
                event_date = event_date.to_pydatetime()

            if isinstance(event_date, datetime):
                raw_round = weekend.get('RoundNumber')
                round_number = int(raw_round) if raw_round is not None and pd.notna(
                    raw_round) else 0

                race_weekends.append(RaceWeekend(
                    round=round_number,
                    full_name=str(weekend.get('OfficialEventName', "") or ""),
                    name=str(weekend.get('EventName', "") or ""),
                    location=str(weekend.get('Location', "") or ""),
                    country=str(weekend.get('Country', "") or ""),
                    start_date=event_date,
                    sessions=sessions
                ))

        return race_weekends

    def get_session_results(self, year: int, round_number: int, session_type: str) -> Optional[SessionResult]:
        try:
            session = fastf1.get_session(year, round_number, session_type)
            session.load(laps=True, telemetry=False,
                         weather=False, messages=False)

            if session.results is None or session.results.empty:
                return None

            winner_row = session.results.loc[session.results['Position'] == 1.0]
            winner_time_delta = None
            if not winner_row.empty:
                winner_time_delta = get_scalar_value(
                    winner_row.iloc[0], 'Time')

            results = []
            for _, row in session.results.iterrows():
                driver_code = str(get_scalar_value(row, 'Abbreviation') or "")
                driver_info = DriverInfo(
                    number=driver_code,
                    full_name=str(get_scalar_value(row, 'FullName') or ""),
                    country_code=str(get_scalar_value(
                        row, 'CountryCode') or ""),
                    team_name=str(get_scalar_value(row, 'TeamName') or "")
                )

                status = str(get_scalar_value(row, 'Status') or "")

                laps_val = get_scalar_value(row, 'Laps')
                if laps_val is None or pd.isna(laps_val):
                    laps_val = get_scalar_value(row, 'NoLaps')

                laps_completed = int(
                    laps_val) if laps_val is not None and pd.notna(laps_val) else 0
                finish_position = int(get_scalar_value(row, 'Position')) if pd.notna(
                    get_scalar_value(row, 'Position')) else 0

                raw_time_delta = get_scalar_value(row, 'Time')
                total_time_ms = None
                gap_ms = None

                normalized_session_type = session_type.upper()

                if normalized_session_type in ['R', 'RACE']:
                    if finish_position == 1:
                        total_time_ms = to_milliseconds(raw_time_delta)
                        gap_ms = 0
                    elif winner_time_delta is not None and raw_time_delta is not None and not pd.isna(raw_time_delta):
                        gap_ms = to_milliseconds(raw_time_delta)
                        total_time_ms = to_milliseconds(
                            winner_time_delta + raw_time_delta)
                else:
                    total_time_ms = to_milliseconds(raw_time_delta)

                fastest_lap_ms = None
                if hasattr(session, 'laps'):
                    driver_laps_any: Any = session.laps.pick_driver(
                        driver_code)
                    if hasattr(driver_laps_any, 'empty') and not driver_laps_any.empty:
                        fastest_lap = driver_laps_any.pick_fastest()
                        if fastest_lap is not None and not pd.isna(fastest_lap['LapTime']):
                            fastest_lap_ms = to_milliseconds(
                                fastest_lap['LapTime'])

                driver_result = DriverResult(
                    position=finish_position,
                    driver=driver_info,
                    laps=laps_completed,
                    status=status,
                    total_time_ms=total_time_ms,
                    gap_ms=gap_ms,
                    fastest_lap_ms=fastest_lap_ms
                )

                if normalized_session_type in ['R', 'RACE']:
                    grid_position = get_scalar_value(row, 'GridPosition')
                    driver_result.race_details = RaceDetails(
                        grid_position=int(grid_position) if pd.notna(
                            grid_position) else 0,
                        status=status,
                        positions_change=0
                    )

                elif normalized_session_type in ['Q', 'QUALIFYING', 'SQ', 'SPRINT QUALIFYING']:
                    driver_result.qualifying_details = QualifyingDetails(
                        q1_ms=to_milliseconds(get_scalar_value(row, 'Q1')),
                        q2_ms=to_milliseconds(get_scalar_value(row, 'Q2')),
                        q3_ms=to_milliseconds(get_scalar_value(row, 'Q3'))
                    )

                results.append(driver_result)

            return SessionResult(
                year=year,
                round=round_number,
                session_type=session_type,
                results=results
            )
        except Exception:
            return None
