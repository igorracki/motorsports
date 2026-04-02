package mappers

import (
	"sort"

	"github.com/igorracki/motorsports/backend/internal/formatters"
	"github.com/igorracki/motorsports/backend/internal/models"
)

func MapRaceWeekend(weekend *models.RaceWeekend) {
	if len(weekend.Sessions) == 0 {
		return
	}

	firstSession := weekend.Sessions[0]
	firstLocalMS := firstSession.TimeUTCMS + firstSession.UTCOffsetMS

	minLocal, maxLocal := firstLocalMS, firstLocalMS
	minUTC, maxUTC := firstSession.TimeUTCMS, firstSession.TimeUTCMS

	for i := range weekend.Sessions {
		session := &weekend.Sessions[i]
		localMS := session.TimeUTCMS + session.UTCOffsetMS

		if localMS < minLocal {
			minLocal = localMS
		}
		if localMS > maxLocal {
			maxLocal = localMS
		}
		if session.TimeUTCMS < minUTC {
			minUTC = session.TimeUTCMS
		}
		if session.TimeUTCMS > maxUTC {
			maxUTC = session.TimeUTCMS
		}

		session.TimeLocal = formatters.FormatTimestamp(localMS)
		session.TimeUTC = formatters.FormatTimestamp(session.TimeUTCMS)
		session.SessionCode = formatters.GetSessionCode(session.Type)
	}

	weekend.StartDateLocalMS = minLocal
	weekend.EndDateLocalMS = maxLocal
	weekend.StartDateUTCMS = minUTC
	weekend.EndDateUTCMS = maxUTC

	weekend.StartDateLocal = formatters.FormatTimestamp(minLocal)
	weekend.EndDateLocal = formatters.FormatTimestamp(maxLocal)
	weekend.StartDateUTC = formatters.FormatTimestamp(minUTC)
	weekend.EndDateUTC = formatters.FormatTimestamp(maxUTC)

	weekend.CountryCode = formatters.GetCountryCode(weekend.Country)
}

func MapSessionResults(results *models.SessionResults) {
	if results == nil || len(results.Results) == 0 {
		return
	}

	isRaceType := isRace(results.SessionType)
	isQualifyingType := isQualifying(results.SessionType)

	allPositionsSet := true
	for _, res := range results.Results {
		if res.Position <= 0 {
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

		if results.Results[i].FastestLapMS == nil {
			return false
		}
		if results.Results[j].FastestLapMS == nil {
			return true
		}
		return *results.Results[i].FastestLapMS < *results.Results[j].FastestLapMS
	})

	if !allPositionsSet || (!isRaceType && !isQualifyingType) {
		for i := range results.Results {
			results.Results[i].Position = i + 1
		}
	}

	var sessionBestLapMS int64 = -1
	if !isRaceType && len(results.Results) > 0 && results.Results[0].FastestLapMS != nil {
		sessionBestLapMS = *results.Results[0].FastestLapMS
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
			result.Gap = "-"
			if gap != 0 {
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
}

func isRace(sessionType string) bool {
	return sessionType == models.SessionTypeRaceShort ||
		sessionType == models.SessionTypeRace ||
		sessionType == models.SessionTypeSprintShort ||
		sessionType == models.SessionTypeSprint
}

func isQualifying(sessionType string) bool {
	return sessionType == models.SessionTypeQualifyingShort ||
		sessionType == models.SessionTypeQualifying ||
		sessionType == models.SessionTypeSprintQualifyingShort ||
		sessionType == models.SessionTypeSprintQualifying
}
