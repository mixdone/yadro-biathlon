package main

import (
	"fmt"

	"github.com/mixdone/yadro-biathlon/config"
	"github.com/mixdone/yadro-biathlon/processor"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Не удалось загрузить конфигурацию:", err)
		return
	}

	events, err := config.LoabEvents("sunny_5_skiers/events")
	if err != nil {
		fmt.Println("Не удалось загрузить события:", err)
		return
	}

	proc := processor.NewProcessor(cfg)
	defer func() {
		proc.FlushLog()
		proc.FlushReport()
	}()

	for _, e := range events {
		proc.ProcessEvent(e)
	}

	proc.PrintFinalReport()
}
