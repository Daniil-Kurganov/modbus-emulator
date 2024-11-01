package main

import (
	"log"
	ta "modbus-emulator/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	var payloads [][]byte
	var err error
	if payloads, err = ta.ParsePackets("coils_read"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
	log.Printf("Payloads: %v", payloads)
}
