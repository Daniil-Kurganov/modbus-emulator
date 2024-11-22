package main

import (
	"log"
	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"
)

func main() {
	log.SetFlags(0)
	var history []structs.HistoryEvent
	var err error
	if history, err = ta.Parser(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	for _, currentHistoryEvent := range history {
		currentHistoryEvent.LogPrint()
	}
}
