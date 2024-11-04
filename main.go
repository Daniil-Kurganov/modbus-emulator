package main

import (
	"log"
	ta "modbus-emulator/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	var history []ta.TCPPacket
	var err error
	if history, err = ta.ParsePackets("coils", "write_32"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
	log.Printf("Payloads: %v", history)
}
