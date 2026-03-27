package services

import (
	"testing"

	"github.com/igorracki/motorsports/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFormatSessionResults_PracticeSorting(t *testing.T) {
	// Given: Unsorted practice results (as seen in the reported issue)
	results := &models.SessionResults{
		SessionType: "FP2",
		Results: []models.DriverResult{
			{
				Position:     1,
				FastestLapMS: int64Ptr(82000), // 1:22.000
				Driver:       models.DriverInfo{FullName: "Lando Norris"},
			},
			{
				Position:     21,
				FastestLapMS: int64Ptr(79729), // 1:19.729 (Should be 1st)
				Driver:       models.DriverInfo{FullName: "Oscar Piastri"},
			},
			{
				Position:     7,
				FastestLapMS: int64Ptr(79943), // 1:19.943 (Should be 2nd)
				Driver:       models.DriverInfo{FullName: "Kimi Antonelli"},
			},
		},
	}

	// When
	formatSessionResults(results)

	// Then: Results should be sorted by FastestLapMS and positions re-assigned
	assert.Equal(t, "Oscar Piastri", results.Results[0].Driver.FullName)
	assert.Equal(t, 1, results.Results[0].Position)
	assert.Equal(t, "1:19.729", results.Results[0].FastestLap)
	assert.Equal(t, "-", results.Results[0].Gap)

	assert.Equal(t, "Kimi Antonelli", results.Results[1].Driver.FullName)
	assert.Equal(t, 2, results.Results[1].Position)
	assert.Equal(t, "+0.214", results.Results[1].Gap)

	assert.Equal(t, "Lando Norris", results.Results[2].Driver.FullName)
	assert.Equal(t, 3, results.Results[2].Position)
	assert.Equal(t, "+2.271", results.Results[2].Gap)
}

func TestFormatSessionResults_RaceOrderPreserved(t *testing.T) {
	// Given: Race results (where position might not be just best lap)
	results := &models.SessionResults{
		SessionType: models.SessionTypeRace,
		Results: []models.DriverResult{
			{
				Position:    1,
				TotalTimeMS: int64Ptr(5400000),
				Driver:      models.DriverInfo{FullName: "Winner"},
			},
			{
				Position:    2,
				TotalTimeMS: int64Ptr(5405000),
				Driver:      models.DriverInfo{FullName: "Second"},
			},
		},
	}

	// When
	formatSessionResults(results)

	// Then: Order should be preserved as per classification from the provider
	assert.Equal(t, "Winner", results.Results[0].Driver.FullName)
	assert.Equal(t, "Second", results.Results[1].Driver.FullName)
}
