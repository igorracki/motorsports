package database

import "strings"

// GeneratePlaceholders returns a string of comma-separated placeholders for a SQL IN clause.
func GeneratePlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat("?,", count-1) + "?"
}

// ToAnySlice converts a slice of any type to a slice of any (interface{}).
// (Useful for passing slices to variadic SQL query arguments).
func ToAnySlice[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
