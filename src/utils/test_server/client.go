package main

import (
	"log"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	var err error
	log.SetFlags(0)
	time.Sleep(time.Second)
	handler := modbus.NewTCPClientHandler("localhost:1502")
	if err = handler.Connect(); err != nil {
		log.Fatalf("Error on handler connecting: %s\n", err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)
	var registers []byte
	for {
		if registers, err = client.ReadCoils(0, 10); err != nil {
			log.Fatalf("Error on reading coils: %s", err)
		}
		log.Printf("\nCoils: %v", registers)
		if registers, err = client.ReadDiscreteInputs(15, 4); err != nil {
			log.Fatalf("Error on reading DI: %s", err)
		}
		log.Printf("DI: %v", registers)
		if registers, err = client.ReadHoldingRegisters(0, 10); err != nil {
			log.Fatalf("Error on reading HR: %s", err)
		}
		log.Printf("HR: %v", registers)
		if registers, err = client.ReadInputRegisters(26, 3); err != nil {
			log.Fatalf("Error on reading IR: %s", err)
		}
		log.Printf("IR`: %v", registers)
	}
}
