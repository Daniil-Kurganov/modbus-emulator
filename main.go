package main

import (
	"log"
	"modbus-emulator/src/server/tcp"
)

func main() {
	log.SetFlags(0)
	tcp.ServerInit()
}
