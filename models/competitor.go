package models

import "time"

type Competitor struct {
	ID   int
	Hits int

	PlannedStart time.Time
	ActualStart  time.Time

	LapStartTimes []time.Time
	LapTimes      []time.Duration

	PenaltyLapStartTimes []time.Time
	PenaltyTime          time.Duration

	Registered  bool
	StartSet    bool
	Started     bool
	Finished    bool
	NotStarted  bool
	NotFinished bool
}

func (c *Competitor) LapTimesSum() time.Duration {
	var total time.Duration
	for _, lap := range c.LapTimes {
		total += lap
	}
	return total
}
