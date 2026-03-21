import fastf1
import numpy
import matplotlib.pyplot as plot


def rotate(matrix, *, angle):
    rotated = numpy.array([
        [numpy.cos(angle), numpy.sin(angle)],
        [-numpy.sin(angle), numpy.cos(angle)]
    ])
    return numpy.matmul(matrix, rotated)


session = fastf1.get_session(2025, 24, 'R')
session.load()

lap = session.laps.pick_fastest()
position = lap.get_pos_data()

circuit_info = session.get_circuit_info()

# Get an array of shape [n, 2] where n is the number of points
# and the second axis is x and y.
track = position.loc[:, ('X', 'Y')].to_numpy()

# Convert the rotation angle from degrees to radian.
track_angle = circuit_info.rotation / 180 * numpy.pi

# Rotate and plot the track map.
rotated_track = rotate(track, angle=track_angle)
plot.plot(rotated_track[:, 0], rotated_track[:, 1])

plot.title(session.event['EventName'])
plot.xticks([])
plot.yticks([])
plot.axis('equal')
plot.show()

print()
print()
print()
