package formatters

import "strings"

var CountryToISO = map[string]string{
	"Bahrain":        "BH",
	"Saudi Arabia":   "SA",
	"Australia":      "AU",
	"Japan":          "JP",
	"China":          "CN",
	"USA":            "US",
	"United States":  "US",
	"Miami":          "US",
	"Monaco":         "MC",
	"Spain":          "ES",
	"Canada":         "CA",
	"Austria":        "AT",
	"United Kingdom": "GB",
	"UK":             "GB",
	"Hungary":        "HU",
	"Belgium":        "BE",
	"Netherlands":    "NL",
	"Italy":          "IT",
	"Azerbaijan":     "AZ",
	"Singapore":      "SG",
	"Mexico":         "MX",
	"Brazil":         "BR",
	"Las Vegas":      "US",
	"Qatar":          "QA",
	"UAE":            "AE",
	"Abu Dhabi":      "AE",
}

var SessionNameToCode = map[string]string{
	"Practice 1":        "FP1",
	"Practice 2":        "FP2",
	"Practice 3":        "FP3",
	"Qualifying":        "Q",
	"Sprint Qualifying": "SQ",
	"Sprint":            "S",
	"Race":              "R",
}

var ISO3ToISO2 = map[string]string{
	"AUS":           "AU",
	"AUT":           "AT",
	"AZE":           "AZ",
	"BEL":           "BE",
	"BRA":           "BR",
	"CAN":           "CA",
	"CHN":           "CN",
	"DEN":           "DK",
	"ESP":           "ES",
	"FIN":           "FI",
	"FRA":           "FR",
	"GBR":           "GB",
	"GER":           "DE",
	"HUN":           "HU",
	"ITA":           "IT",
	"JPN":           "JP",
	"MEX":           "MX",
	"MON":           "MC",
	"NED":           "NL",
	"QAT":           "QA",
	"SGP":           "SG",
	"THA":           "TH",
	"USA":           "US",
	"UAE":           "AE",
	"ARG":           "AR",
	"NZL":           "NZ",
	"MONACO":        "MC",
	"BRITISH":       "GB",
	"DUTCH":         "NL",
	"MEXICAN":       "MX",
	"MONACOESQUE":   "MC",
	"SPANISH":       "ES",
	"AUSTRALIAN":    "AU",
	"CANADIAN":      "CA",
	"FRENCH":        "FR",
	"GERMAN":        "DE",
	"DANISH":        "DK",
	"THAI":          "TH",
	"JAPANESE":      "JP",
	"CHINESE":       "CN",
	"FINNISH":       "FI",
	"AMERICAN":      "US",
	"ITALIAN":       "IT",
	"BRAZILIAN":     "BR",
	"NEW ZEALANDER": "NZ",
	"ARGENTINE":     "AR",
}

var DriverAbbrToISO2 = map[string]string{
	"VER": "NL",
	"PER": "MX",
	"HAM": "GB",
	"RUS": "GB",
	"LEC": "MC",
	"SAI": "ES",
	"NOR": "GB",
	"PIA": "AU",
	"ALO": "ES",
	"STR": "CA",
	"GAS": "FR",
	"OCO": "FR",
	"HUL": "DE",
	"MAG": "DK",
	"ALB": "TH",
	"TSU": "JP",
	"ZHO": "CN",
	"BOT": "FI",
	"SAR": "US",
	"BEA": "GB",
	"ANT": "IT",
	"DOO": "AU",
	"BOR": "BR",
	"HAD": "FR",
	"COL": "AR",
	"LAW": "NZ",
}

func GetCountryCode(country string) string {
	if code, ok := CountryToISO[country]; ok {
		return code
	}
	return ""
}

func GetDriverCountryCode(iso3, abbr string) string {
	if standardized, ok := ISO3ToISO2[strings.ToUpper(iso3)]; ok {
		return standardized
	}
	if standardized, ok := DriverAbbrToISO2[strings.ToUpper(abbr)]; ok {
		return standardized
	}
	return ""
}

func GetSessionCode(sessionName string) string {

	if code, ok := SessionNameToCode[sessionName]; ok {
		return code
	}
	for name, code := range SessionNameToCode {
		if strings.Contains(strings.ToLower(sessionName), strings.ToLower(name)) {
			return code
		}
	}
	return sessionName
}
