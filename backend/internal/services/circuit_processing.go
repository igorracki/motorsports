package services

import (
	"math"

	"github.com/igorracki/motorsports/backend/internal/models"
)

func transformLayout(circuit *models.Circuit) {
	if circuit == nil || len(circuit.Layout) == 0 {
		return
	}

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

func roundCircuitMetrics(circuit *models.Circuit) {
	if circuit == nil || circuit.LengthKm == nil {
		return
	}
	rounded := math.Round(*circuit.LengthKm*1000) / 1000
	circuit.LengthKm = &rounded
}
