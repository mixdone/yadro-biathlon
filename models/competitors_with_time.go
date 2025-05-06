package models

import "time"

type CompetitorWithTime struct {
	Competitor *Competitor
	TotalTime  time.Duration
}

type ByTotalTime []CompetitorWithTime

func (b ByTotalTime) Len() int           { return len(b) }
func (b ByTotalTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByTotalTime) Less(i, j int) bool { return b[i].TotalTime < b[j].TotalTime }
