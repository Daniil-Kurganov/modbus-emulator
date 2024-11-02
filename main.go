package main

import (
	"log"
	ta "modbus-emulator/traffic_analysis"
)

func main() {
	log.SetFlags(0)
	var payloads []ta.TCPPacket
	var err error
	if payloads, err = ta.ParsePackets("coils", "read_01"); err != nil {
		log.Fatalf("Error on parsing file: %v\n", err)
	}
	log.Printf("Payloads: %v", payloads)
}
