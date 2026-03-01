package services

import (
	"github.com/igorracki/f1/backend/internal/formatters"
	"github.com/igorracki/f1/backend/internal/models"
)

func formatSessionResults(results *models.SessionResults) {
	if results == nil || len(results.Results) == 0 {
		return
	}

	for i := range results.Results {
		result := &results.Results[i]

		if result.TotalTimeMS != nil {
			result.TotalTime = formatters.FormatDuration(*result.TotalTimeMS, false)
		}

		if result.FastestLapMS != nil {
			result.FastestLap = formatters.FormatDuration(*result.FastestLapMS, false)
		}

		if results.SessionType == models.SessionTypeRaceShort || results.SessionType == models.SessionTypeRace {
			if result.GapMS != nil {
				result.Gap = formatters.FormatDuration(*result.GapMS, true)
			} else {
				result.Gap = result.Status
			}
		}

		if result.Race != nil {
			result.Race.PositionsChange = result.Race.GridPosition - result.Position
		}

		if result.Qualifying != nil {
			if result.Qualifying.Q1MS != nil {
				result.Qualifying.Q1 = formatters.FormatDuration(*result.Qualifying.Q1MS, false)
			}
			if result.Qualifying.Q2MS != nil {
				result.Qualifying.Q2 = formatters.FormatDuration(*result.Qualifying.Q2MS, false)
			}
			if result.Qualifying.Q3MS != nil {
				result.Qualifying.Q3 = formatters.FormatDuration(*result.Qualifying.Q3MS, false)
			}
		}
	}
}
