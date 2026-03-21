package services

import (
	"github.com/igorracki/f1/backend/internal/formatters"
	"github.com/igorracki/f1/backend/internal/models"
)

func calculateWeekendBoundaries(weekend *models.RaceWeekend) {
	if len(weekend.Sessions) == 0 {
		return
	}

	firstSession := weekend.Sessions[0]
	firstLocalMS := firstSession.TimeUTCMS + firstSession.UTCOffsetMS

	minLocal := firstLocalMS
	maxLocal := firstLocalMS
	minUTC := firstSession.TimeUTCMS
	maxUTC := firstSession.TimeUTCMS

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
	}

	weekend.StartDateLocalMS = minLocal
	weekend.EndDateLocalMS = maxLocal
	weekend.StartDateUTCMS = minUTC
	weekend.EndDateUTCMS = maxUTC
}

func formatRaceWeekend(weekend *models.RaceWeekend) {
	for i := range weekend.Sessions {
		session := &weekend.Sessions[i]
		localMS := session.TimeUTCMS + session.UTCOffsetMS
		session.TimeLocal = formatters.FormatTimestamp(localMS)
		session.TimeUTC = formatters.FormatTimestamp(session.TimeUTCMS)
	}

	weekend.StartDateLocal = formatters.FormatTimestamp(weekend.StartDateLocalMS)
	weekend.EndDateLocal = formatters.FormatTimestamp(weekend.EndDateLocalMS)
	weekend.StartDateUTC = formatters.FormatTimestamp(weekend.StartDateUTCMS)
	weekend.EndDateUTC = formatters.FormatTimestamp(weekend.EndDateUTCMS)
}

func populateStandardCodes(weekend *models.RaceWeekend) {
	weekend.CountryCode = formatters.GetCountryCode(weekend.Country)
	for i := range weekend.Sessions {
		session := &weekend.Sessions[i]
		session.SessionCode = formatters.GetSessionCode(session.Type)
	}
}
