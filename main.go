package main

import (
	"fmt"
	"os"

	"github.com/mixdone/yadro-biathlon/config"
	"github.com/mixdone/yadro-biathlon/processor"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Недостаточно аргументов")
		return
	}

	log, err := os.Create(os.Args[1])
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer log.Close()

	result, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer result.Close()

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

	proc := processor.NewProcessor(cfg, log, result)
	defer func() {
		proc.FlushLog()
		proc.FlushReport()
	}()

	for _, e := range events {
		proc.ProcessEvent(e)
	}

	proc.PrintFinalReport()
}
