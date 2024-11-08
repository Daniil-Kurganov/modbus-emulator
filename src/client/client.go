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
	// values := [][]byte{
	// 	{0, 45, 0, 21},
	// 	{0, 34, 0, 10},
	// }
	// val := []byte{0, 11, 0, 20}
	var registers []byte
	// counter := 0
	for {
		time.Sleep(500 * time.Millisecond)
		if registers, err = client.ReadInputRegisters(4, 1); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		log.Printf("Registers: %v\n", registers)
		// if counter == len(val)-1 {
		// 	counter = 0
		// } else {
		// 	counter += 1
		// }
		// if registers, err = client.WriteMultipleCoils(2, 3, []byte{0, 1, 0}); err != nil {
		// 	log.Fatalf("Error on reading coils: %s\n", err)
		// }
		// log.Printf("Registers: %v\n", registers)
		// time.Sleep(500 * time.Millisecond)
		// if registers, err = client.WriteMultipleRegisters(3, 2, values[0]); err != nil {
		// 	log.Fatalf("Error on writing HR: %s\n", err)
		// }
		// log.Printf("Registers: %v\n", registers)
		// if registers, err = client.WriteMultipleRegisters(3, 2, values[1]); err != nil {
		// 	log.Fatalf("Error on writing HR: %s\n", err)
		// }
		// log.Printf("Registers: %v\n", registers)
	}
}
