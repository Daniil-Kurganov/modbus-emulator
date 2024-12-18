package main

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"time"

	ms "github.com/Daniil-Kurganov/modbus-server"
)

func serverInit(port uint16, protocol string) {
	var err error
	server := ms.NewServer()
	switch protocol {
	case conf.Protocols.RTUOverTCP:
		if err = server.ListenRTUOverTCP(fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
			log.Fatalf("Error on listening RTU over TCP: %s\n", err)
		}
	case conf.Protocols.TCP:
		if err = server.ListenTCP(fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
			log.Fatalf("Error on listening TCP: %s\n", err)
		}
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
	go serverInit(1501, conf.Protocols.RTUOverTCP)
	go serverInit(1502, conf.Protocols.TCP)
	go serverInit(1503, conf.Protocols.RTUOverTCP)
	go serverInit(1504, conf.Protocols.TCP)
	for {
		time.Sleep(time.Second)
	}
}
