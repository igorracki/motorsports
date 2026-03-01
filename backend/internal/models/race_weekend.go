package models

type Session struct {
	Type        string `json:"type"`
	SessionCode string `json:"session_code"`
	TimeLocal   string `json:"time_local,omitempty"`
	TimeUTCMS   int64  `json:"time_utc_ms"`
	TimeUTC     string `json:"time_utc,omitempty"`
	UTCOffsetMS int64  `json:"utc_offset_ms"`
}

type RaceWeekend struct {
	Round            int       `json:"round"`
	FullName         string    `json:"full_name"`
	Name             string    `json:"name"`
	Location         string    `json:"location"`
	Country          string    `json:"country"`
	CountryCode      string    `json:"country_code"`
	EventFormat      string    `json:"event_format"`
	StartDateLocal   string    `json:"start_date_local,omitempty"`
	StartDateLocalMS int64     `json:"start_date_local_ms"`
	StartDateUTC     string    `json:"start_date_utc,omitempty"`
	StartDateUTCMS   int64     `json:"start_date_utc_ms"`
	EndDateLocal     string    `json:"end_date_local,omitempty"`
	EndDateLocalMS   int64     `json:"end_date_local_ms"`
	EndDateUTC       string    `json:"end_date_utc,omitempty"`
	EndDateUTCMS     int64     `json:"end_date_utc_ms"`
	Sessions         []Session `json:"sessions"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type ScheduleResponse struct {
	Schedule []RaceWeekend `json:"schedule"`
}
