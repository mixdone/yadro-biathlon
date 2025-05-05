package main

import (
	"bufio"
)

type Processor struct {
	config      *Config
	competitors map[int]*Competitor
	logWriter   *bufio.Writer
}

func NewProcessor(cfg *Config) *Processor {
	return &Processor{}
}
