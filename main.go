package main

import (
	"log"
	"modbus-emulator/src"
	// ta "modbus-emulator/src/traffic_analysis"
	// "modbus-emulator/src/traffic_analysis/structs"
	// "modbus-emulator/src/utils"
)

func main() {
	log.SetFlags(0)
	src.ServerInit()
	// var history []structs.HistoryEvent
	// var err error
	// if history, err = ta.ParseDump(); err != nil {
	// 	log.Fatalf("Error on parsing dump: %s", err)
	// }
	// log.Printf("Mode: %s", utils.WorkMode)
	// for _, currentHistoryEvent := range history {
	// 	currentHistoryEvent.LogPrint()
	// }
}
