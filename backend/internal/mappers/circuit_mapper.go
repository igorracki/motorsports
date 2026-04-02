package mappers

import (
	"math"

	"github.com/igorracki/motorsports/backend/internal/formatters"
	"github.com/igorracki/motorsports/backend/internal/models"
)

func MapCircuit(circuit *models.Circuit) {
	if circuit == nil {
		return
	}

	circuit.EventDate = formatters.FormatTimestamp(circuit.EventDateMS)

	if circuit.LengthKm != nil {
		rounded := math.Round(*circuit.LengthKm*1000) / 1000
		circuit.LengthKm = &rounded
	}

	if len(circuit.Layout) > 0 {
		rotationRadians := circuit.Rotation * math.Pi / 180.0
		cosineTheta := math.Cos(rotationRadians)
		sineTheta := math.Sin(rotationRadians)

		for i := range circuit.Layout {
			x := circuit.Layout[i].X
			y := circuit.Layout[i].Y

			circuit.Layout[i].X = x*cosineTheta - y*sineTheta
			circuit.Layout[i].Y = -(x*sineTheta + y*cosineTheta)
		}
	}
}
