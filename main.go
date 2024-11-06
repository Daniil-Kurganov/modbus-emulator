package main

import (
	"log"
	ta "modbus-emulator/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	var history map[string]ta.Handshake
	var err error
	if history, err = ta.ParsePackets("HR", "write_32"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
	for currentTransaction, currentHandshake := range history {
		log.Printf("\n\nTransaction № %v\n", currentTransaction)
		log.Println("\n Request:")
		currentHandshake.Request.LogPrint()
		log.Println("\n Responce:")
		currentHandshake.Responce.LogPrint()
	}
}
