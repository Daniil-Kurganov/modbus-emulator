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
	if err = server.ListenTCP(fmt.Sprintf("%s:%s", utils.ServerTCPHost, utils.ServerTCPPort)); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer server.Close()
	server.DiscreteInputs[10], server.DiscreteInputs[13], server.DiscreteInputs[16] = 1, 1, 1
	server.InputRegisters[4], server.InputRegisters[10], server.InputRegisters[18], server.InputRegisters[21] = 120, 385, 16, 6648
	for {
		time.Sleep(500 * time.Millisecond)
	}
}
