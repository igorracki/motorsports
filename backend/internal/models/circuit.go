package models

type CircuitLayoutPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Circuit struct {
	CircuitName string               `json:"circuit_name"`
	Location    string               `json:"location"`
	Country     string               `json:"country"`
	Latitude    float64              `json:"latitude"`
	Longitude   float64              `json:"longitude"`
	LengthKm    float64              `json:"length_km"`
	Corners     int                  `json:"corners"`
	Layout      []CircuitLayoutPoint `json:"layout"`
	EventName   string               `json:"event_name"`
	EventDate   string               `json:"event_date"`
}
