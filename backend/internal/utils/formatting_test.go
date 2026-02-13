package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		ms       int64
		isGap    bool
		expected string
	}{
		{
			name:     "Winner Total Time",
			ms:       5400000, // 1:30:00.000
			isGap:    false,
			expected: "1:30:00.000",
		},
		{
			name:     "Short Gap",
			ms:       12345, // 12.345s
			isGap:    true,
			expected: "+12.345",
		},
		{
			name:     "Long Gap",
			ms:       65432, // 1:05.432
			isGap:    true,
			expected: "+1:05.432",
		},
		{
			name:     "Qualifying Time",
			ms:       90123, // 1:30.123
			isGap:    false,
			expected: "1:30.123",
		},
		{
			name:     "Very Short Gap",
			ms:       500, // 0.500s
			isGap:    true,
			expected: "+0.500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.ms, tt.isGap)
			assert.Equal(t, tt.expected, result)
		})
	}
}
