package services

import (
	"log/slog"
	"sort"

	"github.com/igorracki/f1/backend/internal/formatters"
	"github.com/igorracki/f1/backend/internal/models"
)

func formatSessionResults(results *models.SessionResults) {
	if results == nil || len(results.Results) == 0 {
		return
	}

	slog.Info("Entry: formatSessionResults", "session_type", results.SessionType, "year", results.Year, "round", results.Round, "count", len(results.Results))

	isRaceType := results.SessionType == models.SessionTypeRaceShort ||
		results.SessionType == models.SessionTypeRace ||
		results.SessionType == models.SessionTypeSprintShort ||
		results.SessionType == models.SessionTypeSprint

	isQualifyingType := results.SessionType == models.SessionTypeQualifyingShort ||
		results.SessionType == models.SessionTypeQualifying ||
		results.SessionType == models.SessionTypeSprintQualifyingShort ||
		results.SessionType == models.SessionTypeSprintQualifying

	// Ensure results are sorted correctly.
	// For Practice sessions, we always sort by FastestLapMS.
	// For Race and Qualifying, we trust the Position column if it's populated for all drivers.
	// If any positions are missing (0), we fall back to FastestLapMS for sorting in non-race sessions.
	allPositionsSet := true
	for _, r := range results.Results {
		if r.Position <= 0 {
			allPositionsSet = false
			break
		}
	}

	shouldSortByLapTime := !isRaceType && (!allPositionsSet || !isQualifyingType)

	sort.Slice(results.Results, func(i, j int) bool {
		if !shouldSortByLapTime {
			p1, p2 := results.Results[i].Position, results.Results[j].Position
			if p1 > 0 && p2 > 0 {
				return p1 < p2
			}
			if p1 > 0 {
				return true
			}
			if p2 > 0 {
				return false
			}
		}

		// Fallback or explicit lap time sort
		if results.Results[i].FastestLapMS == nil {
			return false
		}
		if results.Results[j].FastestLapMS == nil {
			return true
		}
		return *results.Results[i].FastestLapMS < *results.Results[j].FastestLapMS
	})

	// Re-assign positions if they were missing or if it's a practice session
	if !allPositionsSet || (!isRaceType && !isQualifyingType) {
		for i := range results.Results {
			results.Results[i].Position = i + 1
		}
	}

	// Find the reference lap for gap calculations in non-race sessions.
	// After potential sorting, the first driver is always the leader.
	var sessionBestLapMS int64 = -1
	if !isRaceType && len(results.Results) > 0 {
		if results.Results[0].FastestLapMS != nil {
			sessionBestLapMS = *results.Results[0].FastestLapMS
		}
	}

	for i := range results.Results {
		result := &results.Results[i]

		result.Driver.CountryCode = formatters.GetDriverCountryCode(result.Driver.CountryCode, result.Driver.ID)

		if result.TotalTimeMS != nil {
			result.TotalTime = formatters.FormatDuration(*result.TotalTimeMS, false)
		}

		if result.FastestLapMS != nil {
			result.FastestLap = formatters.FormatDuration(*result.FastestLapMS, false)
		}

		if isRaceType {
			if result.GapMS != nil {
				result.Gap = formatters.FormatDuration(*result.GapMS, true)
			} else {
				result.Gap = result.Status
			}
		} else if sessionBestLapMS >= 0 && result.FastestLapMS != nil {
			gap := *result.FastestLapMS - sessionBestLapMS
			if gap == 0 {
				result.Gap = "-"
			} else {
				result.Gap = formatters.FormatDuration(gap, true)
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

	slog.Info("Exit: formatSessionResults", "session_type", results.SessionType, "year", results.Year, "round", results.Round, "processed_count", len(results.Results))
}
