package main

import (
	"log"
	// ta "modbus-emulator/src/traffic_analysis"
	s "modbus-emulator/src/server"
)

func main() {
	log.SetFlags(0)
	// var history []ta.History
	// var err error
	// if history, err = ta.ParsePackets("test_files", "HR", "write_42"); err != nil {
	// 	log.Fatalf("Error on parsing file: %v\n", err)
	// }
	// history[1].Print()
	s.Server()
}