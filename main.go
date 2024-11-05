package main

import (
	"log"
	ta "modbus-emulator/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	var history map[string]ta.Handshake
	var err error
	if history, err = ta.ParsePackets("IR", "read_36"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
	for currentKey, currentValue := range history {
		log.Printf("№ %v: %v\n", currentKey, currentValue)
	}
}
