package main

import (
	"fmt"
	"log"
	"modbus-emulator/src/utils"
	"time"

	ms "github.com/Daniil-Kurganov/modbus-server"
)

func serverInit(port uint16) {
	var err error
	server := ms.NewServer()
	if err = server.ListenRTUOverTCP(fmt.Sprintf("%s:%d", utils.ServerTCPHost, port)); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer server.Close()
	server.InitSlave(1)
	server.InitSlave(2)
	server.InitSlave(3)
	log.Print(<-server.ConnectionChanel)
	server.Slaves[1].DiscreteInputs[10], server.Slaves[1].DiscreteInputs[13], server.Slaves[1].DiscreteInputs[16] = 1, 1, 1
	server.Slaves[1].InputRegisters[4], server.Slaves[1].InputRegisters[10], server.Slaves[1].InputRegisters[18], server.Slaves[1].InputRegisters[21] = 120, 385, 16, 6648
	for {
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	log.SetFlags(0)
	go serverInit(1502)
	go serverInit(1503)
	for {
		time.Sleep(time.Second)
	}
}
