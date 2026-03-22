package models

type DriverInfo struct {
	ID          string `json:"id"`
	Number      string `json:"number"`
	FullName    string `json:"full_name"`
	CountryCode string `json:"country_code"`
	TeamName    string `json:"team_name"`
}

type RaceDetails struct {
	GridPosition    int    `json:"grid_position"`
	Status          string `json:"status"`
	PositionsChange int    `json:"positions_change"`
}

type QualifyingDetails struct {
	Q1MS *int64 `json:"q1_ms,omitempty"`
	Q1   string `json:"q1,omitempty"`
	Q2MS *int64 `json:"q2_ms,omitempty"`
	Q2   string `json:"q2,omitempty"`
	Q3MS *int64 `json:"q3_ms,omitempty"`
	Q3   string `json:"q3,omitempty"`
}

type DriverResult struct {
	Position     int                `json:"position"`
	Driver       DriverInfo         `json:"driver"`
	Laps         int                `json:"laps"`
	Status       string             `json:"status"`
	TotalTimeMS  *int64             `json:"total_time_ms,omitempty"`
	TotalTime    string             `json:"total_time,omitempty"`
	GapMS        *int64             `json:"gap_ms,omitempty"`
	Gap          string             `json:"gap,omitempty"`
	FastestLapMS *int64             `json:"fastest_lap_ms,omitempty"`
	FastestLap   string             `json:"fastest_lap,omitempty"`
	Race         *RaceDetails       `json:"race_details,omitempty"`
	Qualifying   *QualifyingDetails `json:"qualifying_details,omitempty"`
}

type SessionResults struct {
	Year        int            `json:"year"`
	Round       int            `json:"round"`
	SessionType string         `json:"session_type"`
	Results     []DriverResult `json:"results"`
}
