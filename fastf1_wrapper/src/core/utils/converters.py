import pandas


def format_lap_time(lap_time) -> str:
    if pandas.isna(lap_time):
        return None

    total_seconds = lap_time.total_seconds()
    minutes = int(total_seconds // 60)
    seconds = int(total_seconds % 60)
    milliseconds = int((total_seconds % 1) * 1000)

    return f"{minutes}:{seconds:02d}.{milliseconds:03d}"
