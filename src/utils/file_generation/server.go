package main

import (
	"fmt"
	"log"
	"modbus-emulator/src/utils"
	"time"

	ms "github.com/Daniil-Kurganov/modbus-server"
)

func main() {
	log.SetFlags(0)
	var err error
	server := ms.NewServer()
	if err = server.ListenRTUOverTCP(fmt.Sprintf("%s:%s", utils.ServerTCPHost, utils.ServerTCPPort)); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer server.Close()
	server.Coils[5], server.HoldingRegisters[4], server.InputRegisters[28] = 1, 16, 103
	for {
		// log.Print("I'm work!")
		time.Sleep(500 * time.Millisecond)
	}
}
