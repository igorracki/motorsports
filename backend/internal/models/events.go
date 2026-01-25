package models

type Session struct {
	Type      string `json:"type"`
	TimeLocal string `json:"time_local"`
	TimeUTC   string `json:"time_utc"`
}

type Event struct {
	Round     int       `json:"round"`
	FullName  string    `json:"full_name"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Country   string    `json:"country"`
	StartDate string    `json:"start_date"`
	Sessions  []Session `json:"sessions"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type EventsResponse struct {
	Events []Event `json:"events"`
}
