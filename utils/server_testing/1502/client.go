package main

import (
	"log"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {
	var err error
	log.SetFlags(0)
	time.Sleep(time.Second)
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
	var coils, DI []bool
	var HR, IR []uint16
	for _, currentUnitID := range []uint8{10} {
		client.SetUnitId(currentUnitID)
		log.Printf("\nCurrent slave: %d", currentUnitID)
		if coils, err = client.ReadCoils(5, 5); err != nil {
			log.Fatalf("Error on reading coils[5:10]: %s", err)
		}
		log.Printf("Coils[5:10] = %v", coils)
		if DI, err = client.ReadDiscreteInputs(9, 11); err != nil {
			log.Fatalf("Error on reading DI[9:20]: %s", err)
		}
		log.Printf("DI[9:20] = %v", DI)
		if HR, err = client.ReadRegisters(110, 15, modbus.HOLDING_REGISTER); err != nil {
			log.Fatalf("Error on read HR[110:15]: %s", err)
		}
		log.Printf("HR[110:15] = %v", HR)
		if IR, err = client.ReadRegisters(4, 18, modbus.INPUT_REGISTER); err != nil {
			log.Fatalf("Error on read IR[4:22]: %s", err)
		}
		log.Printf("IR[4:22] = %v", IR)
	}
	// handler := mc.NewTCPClientHandler("localhost:1507")
	// if err = handler.Connect(); err != nil {
	// 	log.Fatalf("Error on handler connecting: %s\n", err)
	// }
	// handler.SlaveId = 26
	// log.Print("Successfully connect")
	// defer handler.Close()
	// client := mc.NewClient(handler)
	// var registers []byte
	// time.Sleep(2700 * time.Millisecond)
	// if registers, err = client.ReadCoils(0, 10); err != nil {
	// 	log.Fatalf("Error on reading coils: %s", err)
	// }
	// log.Printf("\nCoils: %v", registers)
	// if registers, err = client.ReadDiscreteInputs(20, 10); err != nil {
	// 	log.Fatalf("Error on reading DI: %s", err)
	// }
	// log.Printf("DI: %v", registers)
	// if registers, err = client.ReadHoldingRegisters(0, 10); err != nil {
	// 	log.Fatalf("Error on reading HR: %s", err)
	// }
	// log.Printf("HR: %v", registers)
	// if registers, err = client.ReadInputRegisters(26, 3); err != nil {
	// 	log.Fatalf("Error on reading IR: %s", err)
	// }
	// log.Printf("IR`: %v", registers)
}
