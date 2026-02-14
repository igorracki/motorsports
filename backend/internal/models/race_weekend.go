package models

type Session struct {
	Type        string `json:"type"`
	TimeLocalMS int64  `json:"time_local_ms"`
	TimeLocal   string `json:"time_local,omitempty"`
	TimeUTCMS   int64  `json:"time_utc_ms"`
	TimeUTC     string `json:"time_utc,omitempty"`
}

type RaceWeekend struct {
	Round       int       `json:"round"`
	FullName    string    `json:"full_name"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Country     string    `json:"country"`
	StartDateMS int64     `json:"start_date_ms"`
	StartDate   string    `json:"start_date,omitempty"`
	Sessions    []Session `json:"sessions"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type ScheduleResponse struct {
	Schedule []RaceWeekend `json:"schedule"`
}
