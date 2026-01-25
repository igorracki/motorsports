import matplotlib.pyplot as plot
import fastf1.plotting

# Load FastF1's dark color scheme
fastf1.plotting.setup_mpl(mpl_timedelta_support=False, color_scheme='fastf1')

session = fastf1.get_session(2025, 24, 'R')
session.load(telemetry=False, weather=False)

figure, axis = plot.subplots(figsize=(8.0, 4.9))

for driver in session.drivers:
    driver_laps = session.laps.pick_drivers(driver)

    abbreviation = driver_laps['Driver'].iloc[0]
    style = fastf1.plotting.get_driver_style(identifier=abbreviation,
                                             style=['color', 'linestyle'],
                                             session=session)

    axis.plot(driver_laps['LapNumber'],
              driver_laps['Position'], label=abbreviation, **style)

axis.set_ylim([20.5, 0.5])
axis.set_yticks([1, 5, 10, 15, 20])
axis.set_xlabel('Lap')
axis.set_ylabel('Position')

axis.legend(bbox_to_anchor=(1.0, 1.02))
plot.tight_layout()
plot.show()
