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
