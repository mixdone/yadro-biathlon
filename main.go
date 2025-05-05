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

	processor := processor.NewProcessor(cfg)
	defer processor.FlushLog()

	for _, e := range events {
		processor.ProcessEvent(e)
	}

}
