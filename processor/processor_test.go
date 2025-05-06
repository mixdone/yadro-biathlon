package processor

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/mixdone/yadro-biathlon/config"
	"github.com/mixdone/yadro-biathlon/models"
	"github.com/stretchr/testify/assert"
)

func TestProcessor_ProcessEvent_RegisterCompetitor(t *testing.T) {
	cfg := &config.Config{
		Laps:        3,
		LapLen:      5000,
		FiringLines: 1,
	}
	logBuffer := &bytes.Buffer{}
	resultBuffer := &bytes.Buffer{}
	processor := NewProcessor(cfg, logBuffer, resultBuffer)

	event := config.Event{
		EventID:      1,
		CompetitorID: 1,
		Time:         time.Now(),
	}

	processor.ProcessEvent(event)
	processor.FlushLog()

	expectedOutput := fmt.Sprintf("[%s] The competitor(1) registered\n", event.Time.Format("15:04:05.000"))

	assert.Contains(t, logBuffer.String(), expectedOutput)
}

func TestProcessor_PrintFinalReport_SortedByTotalTime(t *testing.T) {
	cfg := &config.Config{
		Laps:        3,
		LapLen:      600,
		FiringLines: 1,
	}
	logBuffer := &bytes.Buffer{}
	resultBuffer := &bytes.Buffer{}
	processor := NewProcessor(cfg, logBuffer, resultBuffer)

	competitor1 := models.Competitor{
		ID:       1,
		Started:  true,
		LapTimes: []time.Duration{1 * time.Minute, 1 * time.Minute, 1 * time.Minute},
		Hits:     5,
	}
	competitor2 := models.Competitor{
		ID:       2,
		Started:  true,
		LapTimes: []time.Duration{2 * time.Minute, 2 * time.Minute, 2 * time.Minute},
		Hits:     5,
	}

	processor.competitors[1] = &competitor1
	processor.competitors[2] = &competitor2

	processor.PrintFinalReport()
	processor.FlushReport()

	expectedOutput1 := "[00:03:00.000] 1 [{00:01:00.000, 10.000}, {00:01:00.000, 10.000}, {00:01:00.000, 10.000}] {00:00:00.000, 0.000} 5/5\n"
	expectedOutput2 := "[00:06:00.000] 2 [{00:02:00.000, 5.000}, {00:02:00.000, 5.000}, {00:02:00.000, 5.000}] {00:00:00.000, 0.000} 5/5\n"

	assert.Contains(t, resultBuffer.String(), expectedOutput1)
	assert.Contains(t, resultBuffer.String(), expectedOutput2)
}

func TestProcessor_LapTimesSum(t *testing.T) {
	competitor := models.Competitor{ID: 1, Started: true}
	competitor.LapTimes = append(competitor.LapTimes, 2*time.Minute)
	competitor.LapTimes = append(competitor.LapTimes, 3*time.Minute)
	competitor.LapTimes = append(competitor.LapTimes, 4*time.Minute)

	expectedSum := 9 * time.Minute
	actualSum := competitor.LapTimesSum()
	assert.Equal(t, expectedSum, actualSum)
}

func TestFormatDuration(t *testing.T) {
	d := 2*time.Hour + 3*time.Minute + 45*time.Second + 123*time.Millisecond
	formatted := formatDuration(d)
	expected := "02:03:45.123"
	assert.Equal(t, expected, formatted)
}
