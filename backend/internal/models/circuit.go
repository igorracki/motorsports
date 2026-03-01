package models

type CircuitLayoutPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Circuit struct {
	CircuitName  string               `json:"circuit_name"`
	Location     string               `json:"location"`
	Country      string               `json:"country"`
	Latitude     *float64             `json:"latitude,omitempty"`
	Longitude    *float64             `json:"longitude,omitempty"`
	LengthKm     *float64             `json:"length_km,omitempty"`
	Corners      *int                 `json:"corners,omitempty"`
	Layout       []CircuitLayoutPoint `json:"layout"`
	EventName    string               `json:"event_name"`
	EventDateMS  int64                `json:"event_date_ms"`
	EventDate    string               `json:"event_date,omitempty"`
	Rotation     float64              `json:"rotation"`
	MaxSpeedKmh  float64              `json:"max_speed_kmh"`
	MaxAltitudeM float64              `json:"max_altitude_m"`
	MinAltitudeM float64              `json:"min_altitude_m"`
}
