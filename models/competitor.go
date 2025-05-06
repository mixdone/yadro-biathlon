package models

import "time"

type Competitor struct {
	ID            int
	Registered    bool
	PlannedStart  time.Time
	ActualStart   time.Time
	StartSet      bool
	Started       bool
	Finished      bool
	NotStarted    bool
	NotFinished   bool
	LapTimes      []time.Duration
	PenaltyTime   time.Duration
	ShotsFired    int
	Hits          int
	LapStartTimes []time.Time
}

func (c *Competitor) LapTimesSum() time.Duration {
	total := time.Duration(0)
	for _, lap := range c.LapTimes {
		total += lap
	}
	return total + c.PenaltyTime
}
