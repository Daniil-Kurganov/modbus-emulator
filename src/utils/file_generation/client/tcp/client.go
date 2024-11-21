package main

import (
	"log"
	"time"

	mc "github.com/goburrow/modbus"
)

func main() {
	var err error
	log.SetFlags(0)
	handler := mc.NewTCPClientHandler("localhost:1502")
	if err = handler.Connect(); err != nil {
		log.Fatalf("Error on handler connecting: %s\n", err)
	}
	defer handler.Close()
	client := mc.NewClient(handler)
	var registers []byte
	if registers, err = client.ReadCoils(0, 1); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read coils (0, 1): %v\n", registers)
	time.Sleep(500 * time.Millisecond)
	if registers, err = client.ReadDiscreteInputs(16, 2); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read DI (16, 2): %v\n", registers)
	time.Sleep(901 * time.Millisecond)
	if registers, err = client.WriteSingleRegister(8, 39); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Write HR[8] = 39: %v\n", registers)
	time.Sleep(1 * time.Second)
	if registers, err = client.ReadCoils(5, 5); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read coils (5, 5): %v\n", registers)
	time.Sleep(100 * time.Millisecond)
	if registers, err = client.WriteMultipleCoils(4, 4, []byte{1, 1, 0, 1}); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Write coils[4:8] = {1, 1, 0, 1}: %v\n", registers)
	time.Sleep(1020 * time.Millisecond)
	if registers, err = client.ReadInputRegisters(25, 4); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read IR (25, 4): %v\n", registers)
}
