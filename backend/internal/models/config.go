package models

type DriverMetadata struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	TeamName    string `json:"team_name"`
	TeamColor   string `json:"team_color"`
	CountryCode string `json:"country_code"`
}

type ValidationConfig struct {
	MinYear    int `json:"min_year"`
	MaxYear    int `json:"max_year"`
	MinRound   int `json:"min_round"`
	MaxRound   int `json:"max_round"`
	MinEntries int `json:"min_entries"`
	MaxEntries int `json:"max_entries"`
}

type AppConfig struct {
	Drivers         []DriverMetadata  `json:"drivers"`
	SessionMappings map[string]string `json:"session_mappings"`
	Validation      ValidationConfig  `json:"validation"`
}
