package main

import (
	"log"
	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"
	"modbus-emulator/src/utils"
)

func main() {
	log.SetFlags(0)
	var history []structs.HistoryEvent
	var err error
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	log.Printf("Mode: %s", utils.Mode)
	for _, currentHistoryEvent := range history {
		currentHistoryEvent.LogPrint()
	}
}
