package utils

import "fmt"

func FormatDuration(ms int64, isGap bool) string {
	prefix := ""
	if isGap {
		prefix = "+"
	}

	secondsTotal := float64(ms) / 1000.0

	if secondsTotal < 60 && isGap {
		return fmt.Sprintf("%s%.3f", prefix, secondsTotal)
	}

	hours := int(secondsTotal) / 3600
	minutes := (int(secondsTotal) % 3600) / 60
	seconds := secondsTotal - float64(hours*3600+minutes*60)

	if hours > 0 {
		return fmt.Sprintf("%s%d:%02d:%06.3f", prefix, hours, minutes, seconds)
	}

	return fmt.Sprintf("%s%d:%06.3f", prefix, minutes, seconds)
}
