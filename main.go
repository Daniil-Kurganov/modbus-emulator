package main

import (
	"log"
	ta "modbus-emulator/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	if err := ta.ParsePackets("coils_read"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
}
