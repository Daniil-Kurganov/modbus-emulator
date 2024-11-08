package main

import (
	"log"
	ta "modbus-emulator/src/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	var history []ta.History
	var err error
	if history, err = ta.ParsePackets("test_files", "IR", "read_41"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
	for _, currentHistoryEvent := range history {
		log.Printf("\n\nTransaction № %v\n", currentHistoryEvent.TransactionID)
		log.Println("\n Request:")
		currentHistoryEvent.Handshake.Request.LogPrint()
		log.Println("\n Response:")
		currentHistoryEvent.Handshake.Response.LogPrint()
	}
}
