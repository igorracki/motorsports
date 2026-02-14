package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name       string
		durationMS int64
		isGap      bool
		expected   string
	}{
		{
			name:       "Winner Total Time",
			durationMS: 5400000,
			isGap:      false,
			expected:   "1:30:00.000",
		},
		{
			name:       "Short Gap",
			durationMS: 12345,
			isGap:      true,
			expected:   "+12.345",
		},
		{
			name:       "Long Gap",
			durationMS: 65432,
			isGap:      true,
			expected:   "+1:05.432",
		},
		{
			name:       "Qualifying Time",
			durationMS: 90123,
			isGap:      false,
			expected:   "1:30.123",
		},
		{
			name:       "Very Short Gap",
			durationMS: 500,
			isGap:      true,
			expected:   "+0.500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.durationMS, tt.isGap)
			assert.Equal(t, tt.expected, result)
		})
	}
}
