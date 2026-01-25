#!/usr/bin/env python3

import argparse
from core.providers import FastF1Provider
from core.services import SessionService, EventService


def main():
    parser = argparse.ArgumentParser(description="FastF1 Wrapper")
    parser.add_argument(
        '--year',
        required=True,
        type=int,
        help='Season year'
    )
    parser.add_argument(
        '--round',
        required=False,
        type=int,
        help='Round number'
    )
    parser.add_argument(
        '--session',
        type=str,
        default='R',
        help='Session type (R,Q,FP1)'
    )
    arguments = parser.parse_args()

    provider = FastF1Provider()
    # session_service = SessionService(provider=provider)
    # race_result = session_service.get_race_result(
    #    arguments.year,
    #    arguments.round,
    #    arguments.session
    # )
    # print(race_result.json())

    event_service = EventService(provider=provider)
    events = event_service.get_weekend_events(arguments.year)
    for event in events:
        print(event.json())


if __name__ == '__main__':
    main()
