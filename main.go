package main

import (
	"log"
	// ta "modbus-emulator/src/traffic_analysis"
	s "modbus-emulator/src/server"
)

func main() {
	log.SetFlags(0)
	s.Server()
}
