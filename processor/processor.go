package processor

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mixdone/yadro-biathlon/config"
	"github.com/mixdone/yadro-biathlon/models"
)

type Processor struct {
	config      *config.Config
	competitors map[int]*models.Competitor
	logWriter   *bufio.Writer
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		config:      cfg,
		competitors: make(map[int]*models.Competitor),
		logWriter:   bufio.NewWriter(os.Stdout),
	}
}

func (p *Processor) logEvent(t time.Time, msg string) {
	fmt.Fprintf(p.logWriter, "[%s] %s\n", t.Format("15:04:05.000"), msg)
}

func (p *Processor) FlushLog() {
	p.logWriter.Flush()
}

func (p *Processor) getOrCreateComp(id int) *models.Competitor {
	if _, ok := p.competitors[id]; !ok {
		p.competitors[id] = &models.Competitor{ID: id}
	}
	return p.competitors[id]
}

func (p *Processor) ProcessEvent(e config.Event) {
	c := p.getOrCreateComp(e.CompetitorID)

	switch e.EventID {
	case 1:
		c.Registered = true
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) registered", c.ID))

	case 2:
		t, _ := time.Parse("15:04:05.000", e.Extra)
		c.PlannedStart = t
		c.StartSet = true
		p.logEvent(e.Time, fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %v",
			c.ID, c.PlannedStart.Format("15:04:05.000")))

	case 3:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) is on the start line", c.ID))

	case 4:
		c.Started = true
		c.ActualStart = e.Time
		c.LapStartTimes = append(c.LapStartTimes, e.Time)
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) has started", c.ID))

	case 5:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) is on the firing range(%s)", c.ID, e.Extra))

	case 6:
		c.ShotsFired++
		c.Hits++
		p.logEvent(e.Time, fmt.Sprintf("The target(%s) has been hit by competitor(%d)", e.Extra, c.ID))

	case 7:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) left the firing range", c.ID))

	case 8:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) entered the penalty laps", c.ID))
		c.LapStartTimes = append(c.LapStartTimes, e.Time)

	case 9:
		penaltyStart := c.LapStartTimes[len(c.LapStartTimes)-1]
		c.PenaltyTime += e.Time.Sub(penaltyStart)
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) left the penalty laps", c.ID))

	case 10:
		lastLapStart := c.LapStartTimes[len(c.LapStartTimes)-1]
		c.LapTimes = append(c.LapTimes, e.Time.Sub(lastLapStart))
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) ended the main lap", c.ID))

	case 11:
		c.NotFinished = true
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) can`t continue: %s", c.ID, e.Extra))
	}
}
