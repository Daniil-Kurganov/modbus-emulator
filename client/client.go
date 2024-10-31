package main

import (
	"log"
	"time"

	"github.com/goburrow/modbus"
)

func main() {
	var err error
	log.SetFlags(0)
	handler := modbus.NewTCPClientHandler("localhost:1502")
	if err = handler.Connect(); err != nil {
		log.Fatalf("Error on handler connecting: %s\n", err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)
	var registers []byte
	for {
		time.Sleep(500 * time.Millisecond)
		if registers, err = client.ReadInputRegisters(0, 1); err != nil {
			log.Fatalf("Error on reading input registers: %s\n", err)
		}
		log.Printf("Registers: %v\n", registers)
	}
}
