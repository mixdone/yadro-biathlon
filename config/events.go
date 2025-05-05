package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	Extra        string
	RawLine      string
}

func parseLine(line string) (Event, error) {
	parts := strings.SplitN(line, "]", 2)
	timeStr := strings.TrimPrefix(parts[0], "[")
	eventTime, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return Event{}, fmt.Errorf("Unable to parse time, %s", err)
	}

	part1 := strings.TrimPrefix(parts[1], " ")
	fields := strings.SplitN(part1, " ", 3)
	eventId, _ := strconv.Atoi(fields[0])
	competitorId, _ := strconv.Atoi(fields[1])
	extra := ""
	if len(fields) == 3 {
		extra = fields[2]
	}

	return Event{
		Time:         eventTime,
		EventID:      eventId,
		CompetitorID: competitorId,
		Extra:        extra,
		RawLine:      line}, nil
}

func LoabEvents(filename string) ([]Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Can't open file %s", err)
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		event, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse line, %s", err)
		}
		events = append(events, event)
	}

	return events, nil
}
