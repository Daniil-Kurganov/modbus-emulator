package main

import (
	"log"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	log.SetFlags(0)
	var err error
	var client *modbus.ModbusClient
	if client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "rtuovertcp://localhost:1502",
		Speed:   19200,
		Timeout: 1 * time.Second,
	}); err != nil {
		log.Fatalf("Error on creating client: %s", err)
	}
	if err = client.Open(); err != nil {
		log.Fatalf("Error on openning client connection: %s", err)
	}
	defer client.Close()
	var coil0 bool
	if coil0, err = client.ReadCoil(0); err != nil {
		log.Fatalf("Error on getting coils[0]: %s", err)
	}
	log.Printf("Read coils (0, 1): %v\n", coil0)
	time.Sleep(500 * time.Millisecond)
	var DI1617 []bool
	if DI1617, err = client.ReadDiscreteInputs(16, 2); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read DI (16, 2): %v\n", DI1617)
	time.Sleep(901 * time.Millisecond)
	if err = client.WriteRegister(8, 39); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Print("Write HR[8] = 39: success")
	time.Sleep(1 * time.Second)
	var coils59 []bool
	if coils59, err = client.ReadCoils(5, 5); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read coils[5:10]: %v\n", coils59)
	time.Sleep(100 * time.Millisecond)
	if err = client.WriteCoils(4, []bool{true, true, false, true}); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Print("Write coils[4:8] = {1, 1, 0, 1}: success")
	time.Sleep(1020 * time.Millisecond)
	var HR37 []uint16
	if HR37, err = client.ReadRegisters(3, 4, modbus.HOLDING_REGISTER); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	log.Printf("Read HR[3:8]: %v\n", HR37)
}
