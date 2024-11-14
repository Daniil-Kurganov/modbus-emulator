package main

import (
	"fmt"
	"log"
	"modbus-emulator/src/utils"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
	var err error
	server := mbserver.NewServer()
	if err = server.ListenTCP(fmt.Sprintf("%s:%s", utils.ServerTCPHost, utils.ServerTCPPort)); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer server.Close()
	server.Coils[5], server.HoldingRegisters[4], server.InputRegisters[28] = 1, 16, 103
	for {
		time.Sleep(500 * time.Millisecond)
	}
}
