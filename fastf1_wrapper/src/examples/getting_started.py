import fastf1

session = fastf1.get_session(2025, 'Abu Dhabi', 'Q')

print(session.name)
print()
print(session.date)
print()
print(session.event)
print()
print(session.event['EventName'])
print()
print(session.event['EventDate'])
print()

# event = fastf1.get_event(2025, 23)
event = fastf1.get_event(2025, 'Qatar')
print(event)
print()

schedule = fastf1.get_event_schedule(2026)
print(schedule)
print()

gp_12 = schedule.get_event_by_round(12)
print(gp_12['Country'])
print()
gp_austin = schedule.get_event_by_name('Austin')
print(gp_austin['Country'])
print()

session.load()
print(session.results)
print()

top_10_quali = session.results.iloc[0:10].loc[:, ['Abbreviation', 'Q3']]
print(top_10_quali)
print()

print(session.laps)
print()

fastest_lap = session.laps.pick_fastest()
print(fastest_lap['LapTime'])
print()
print(fastest_lap['Driver'])
print()
