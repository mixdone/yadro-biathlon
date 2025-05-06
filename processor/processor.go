package processor

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/mixdone/yadro-biathlon/config"
	"github.com/mixdone/yadro-biathlon/models"
)

type Processor struct {
	config       *config.Config
	competitors  map[int]*models.Competitor
	logWriter    *bufio.Writer
	resultWriter *bufio.Writer
}

func NewProcessor(cfg *config.Config, log io.Writer, result io.Writer) *Processor {
	return &Processor{
		config:       cfg,
		competitors:  make(map[int]*models.Competitor),
		logWriter:    bufio.NewWriter(log),
		resultWriter: bufio.NewWriter(result),
	}
}

func (p *Processor) logEvent(t time.Time, msg string) {
	fmt.Fprintf(p.logWriter, "[%s] %s\n", t.Format("15:04:05.000"), msg)
}

func (p *Processor) FlushLog() {
	p.logWriter.Flush()
}

func (p *Processor) FlushReport() {
	p.resultWriter.Flush()
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
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) is on the firing range(%s)",
			c.ID, e.Extra))

	case 6:
		c.Hits++
		p.logEvent(e.Time, fmt.Sprintf("The target(%s) has been hit by competitor(%d)",
			e.Extra, c.ID))

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
		mainLapDuration := e.Time.Sub(c.PlannedStart)
		c.LapTimes = append(c.LapTimes, mainLapDuration)
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) ended the main lap", c.ID))

	case 11:
		c.NotFinished = true
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) can`t continue: %s",
			c.ID, e.Extra))
	}
}

func (p *Processor) PrintFinalReport() {
	shotsFired := 5 * p.config.FiringLines

	var competitorsWithTime []models.CompetitorWithTime
	for _, c := range p.competitors {
		var totalTime time.Duration
		if !c.NotFinished {
			totalTime = c.LapTimesSum()
		}
		competitorsWithTime = append(competitorsWithTime, models.CompetitorWithTime{
			Competitor: c,
			TotalTime:  totalTime,
		})
	}

	sort.Sort(models.ByTotalTime(competitorsWithTime))

	for _, cwt := range competitorsWithTime {
		c := cwt.Competitor
		if !c.Started && c.StartSet {
			fmt.Fprintf(p.resultWriter, "[NotStarted] %d\n", c.ID)
			continue
		}
		if c.NotFinished {
			fmt.Fprintf(p.resultWriter, "[NotFinished] %d ", c.ID)
		} else {
			totalTime := c.LapTimesSum()
			fmt.Fprintf(p.resultWriter, "[%s] %d ", formatDuration(totalTime), c.ID)
		}

		fmt.Fprintf(p.resultWriter, "[")
		for i := 0; i < p.config.Laps; i++ {
			if i > 0 {
				fmt.Fprintf(p.resultWriter, ", ")
			}

			if i < len(c.LapTimes) {
				lap := c.LapTimes[i]
				speed := float64(p.config.LapLen) / float64(lap.Seconds())
				fmt.Fprintf(p.resultWriter, "{%s, %.3f}", formatDuration(lap), speed)
			} else {
				fmt.Fprintf(p.resultWriter, "{,}")
			}
		}
		fmt.Fprintf(p.resultWriter, "] ")

		if c.PenaltyTime >= 0 {
			penaltyLaps := shotsFired - c.Hits
			if penaltyLaps > 0 {
				speed := float64(p.config.PenaltyLen*penaltyLaps) / c.PenaltyTime.Seconds()
				fmt.Fprintf(p.resultWriter, "{%s, %.3f} ", formatDuration(c.PenaltyTime), speed)
			} else {
				fmt.Fprintf(p.resultWriter, "{00:00:00.000, 0.000} ")
			}
		} else {
			fmt.Fprintf(p.resultWriter, "{,} ")
		}

		fmt.Fprintf(p.resultWriter, "%d/%d\n", c.Hits, shotsFired)
	}
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
