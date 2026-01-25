from dataclasses import dataclass


@dataclass
class Driver:
    driver_number: int
    first_name: str
    surname: str
    team: str
