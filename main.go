package main

import (
	"log"
	"modbus-emulator/src"
)

func main() {
	log.SetFlags(0)
	src.ServerInit()
}
