package formatters

import (
	"fmt"
	"time"
)

func FormatDuration(durationMS int64, isGap bool) string {
	prefix := ""
	if isGap && durationMS >= 0 {
		prefix = "+"
	}

	secondsTotal := float64(durationMS) / 1000.0

	if secondsTotal < 60 && secondsTotal > -60 && isGap {
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

func FormatTimestamp(ms int64) string {
	timestamp := time.UnixMilli(ms).UTC()
	return timestamp.Format("2006-01-02T15:04:05")
}
