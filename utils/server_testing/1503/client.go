package main

import (
	"log"
	"time"

	mc "github.com/goburrow/modbus"
)

func main() {
	var err error
	log.SetFlags(0)
	time.Sleep(time.Second)
	handler := mc.NewTCPClientHandler("localhost:1503")
	if err = handler.Connect(); err != nil {
		log.Fatalf("Error on handler connecting: %s\n", err)
	}
	handler.SlaveId = 1
	log.Print("Successfully connect")
	defer handler.Close()
	client := mc.NewClient(handler)
	var registers []byte
	time.Sleep(2700 * time.Millisecond)
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
