package services

import (
	"strings"

	"github.com/igorracki/motorsports/backend/internal/formatters"
	"github.com/igorracki/motorsports/backend/internal/models"
)

type ConfigService interface {
	GetAppConfig() models.AppConfig
}

type configService struct {
	drivers         []models.DriverMetadata
	sessionMappings map[string]string
	validation      models.ValidationConfig
}

func NewConfigService() ConfigService {
	drivers := []models.DriverMetadata{
		{ID: "VER", FullName: "M. Verstappen", TeamName: "Red Bull Racing", TeamColor: "#3671C6", CountryCode: "NL"},
		{ID: "PER", FullName: "S. Perez", TeamName: "Red Bull Racing", TeamColor: "#3671C6", CountryCode: "MX"},
		{ID: "NOR", FullName: "L. Norris", TeamName: "McLaren", TeamColor: "#FF8000", CountryCode: "GB"},
		{ID: "PIA", FullName: "O. Piastri", TeamName: "McLaren", TeamColor: "#FF8000", CountryCode: "AU"},
		{ID: "LEC", FullName: "C. Leclerc", TeamName: "Ferrari", TeamColor: "#E80020", CountryCode: "MC"},
		{ID: "HAM", FullName: "L. Hamilton", TeamName: "Ferrari", TeamColor: "#E80020", CountryCode: "GB"},
		{ID: "RUS", FullName: "G. Russell", TeamName: "Mercedes", TeamColor: "#27F4D2", CountryCode: "GB"},
		{ID: "ANT", FullName: "K. Antonelli", TeamName: "Mercedes", TeamColor: "#27F4D2", CountryCode: "IT"},
		{ID: "ALO", FullName: "F. Alonso", TeamName: "Aston Martin", TeamColor: "#229971", CountryCode: "ES"},
		{ID: "STR", FullName: "L. Stroll", TeamName: "Aston Martin", TeamColor: "#229971", CountryCode: "CA"},
		{ID: "GAS", FullName: "P. Gasly", TeamName: "Alpine", TeamColor: "#0093CC", CountryCode: "FR"},
		{ID: "DOO", FullName: "J. Doohan", TeamName: "Alpine", TeamColor: "#0093CC", CountryCode: "AU"},
		{ID: "ALB", FullName: "A. Albon", TeamName: "Williams", TeamColor: "#64C4FF", CountryCode: "TH"},
		{ID: "SAI", FullName: "C. Sainz", TeamName: "Williams", TeamColor: "#64C4FF", CountryCode: "ES"},
		{ID: "TSU", FullName: "Y. Tsunoda", TeamName: "RB", TeamColor: "#6692FF", CountryCode: "JP"},
		{ID: "HAD", FullName: "I. Hadjar", TeamName: "RB", TeamColor: "#6692FF", CountryCode: "FR"},
		{ID: "HUL", FullName: "N. Hulkenberg", TeamName: "Sauber", TeamColor: "#52E252", CountryCode: "DE"},
		{ID: "BOR", FullName: "G. Bortoleto", TeamName: "Sauber", TeamColor: "#52E252", CountryCode: "BR"},
		{ID: "OCO", FullName: "E. Ocon", TeamName: "Haas", TeamColor: "#B6BABD", CountryCode: "FR"},
		{ID: "BEA", FullName: "O. Bearman", TeamName: "Haas", TeamColor: "#B6BABD", CountryCode: "GB"},
		{ID: "BOT", FullName: "V. Bottas", TeamName: "Sauber", TeamColor: "#52E252", CountryCode: "FI"},
		{ID: "ZHO", FullName: "G. Zhou", TeamName: "Sauber", TeamColor: "#52E252", CountryCode: "CN"},
		{ID: "MAG", FullName: "K. Magnussen", TeamName: "Haas", TeamColor: "#B6BABD", CountryCode: "DK"},
		{ID: "SAR", FullName: "L. Sargeant", TeamName: "Williams", TeamColor: "#64C4FF", CountryCode: "US"},
		{ID: "COL", FullName: "F. Colapinto", TeamName: "Williams", TeamColor: "#64C4FF", CountryCode: "AR"},
		{ID: "LAW", FullName: "L. Lawson", TeamName: "RB", TeamColor: "#6692FF", CountryCode: "NZ"},
	}

	for i := range drivers {
		drivers[i].ID = strings.ToUpper(drivers[i].ID)
	}

	return &configService{
		drivers:         drivers,
		sessionMappings: formatters.SessionNameToCode,
		validation: models.ValidationConfig{
			MinYear:    models.MinF1Year,
			MaxYear:    models.MaxF1Year,
			MinRound:   models.MinF1Round,
			MaxRound:   models.MaxF1Round,
			MinEntries: 3,
			MaxEntries: 22,
		},
	}
}

func (s *configService) GetAppConfig() models.AppConfig {
	return models.AppConfig{
		Drivers:         s.drivers,
		SessionMappings: s.sessionMappings,
		Validation:      s.validation,
	}
}
