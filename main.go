package main

import (
	"log"
	// ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src"
)

func main() {
	log.SetFlags(0)
	src.ServerInit()
}
