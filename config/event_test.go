package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		line        string
		expected    Event
		expectedErr bool
	}{
		{
			line: "[09:31:49.285] 1 3",
			expected: Event{
				Time:         time.Date(0, 1, 1, 9, 31, 49, 285000000, time.UTC),
				EventID:      1,
				CompetitorID: 3,
				Extra:        "",
				RawLine:      "[09:31:49.285] 1 3",
			},
			expectedErr: false,
		},
		{
			line: "[09:55:00.000] 2 1 10:00:00.000",
			expected: Event{
				Time:         time.Date(0, 1, 1, 9, 55, 0, 0, time.UTC),
				EventID:      2,
				CompetitorID: 1,
				Extra:        "10:00:00.000",
				RawLine:      "[09:55:00.000] 2 1 10:00:00.000",
			},
			expectedErr: false,
		},
		{
			line:        "[invalid_time] 1 3",
			expected:    Event{},
			expectedErr: true,
		},
		{
			line: "[09:31:49.285] 1 3",
			expected: Event{
				Time:         time.Date(0, 1, 1, 9, 31, 49, 285000000, time.UTC),
				EventID:      1,
				CompetitorID: 3,
				Extra:        "",
				RawLine:      "[09:31:49.285] 1 3",
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result, err := parseLine(tt.line)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
