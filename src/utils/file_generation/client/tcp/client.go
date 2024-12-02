package main

import (
	"log"

	mc "github.com/goburrow/modbus"
)

func main() {
	var err error
	log.SetFlags(0)
	for _, currentSlaveID := range []byte{1, 2, 3} {
		handler := mc.NewTCPClientHandler("localhost:1502")
		if err = handler.Connect(); err != nil {
			log.Fatalf("Error on handler connecting: %s\n", err)
		}
		handler.SlaveId = currentSlaveID
		log.Printf("\nSet slave ID = %d", currentSlaveID)
		client := mc.NewClient(handler)
		var registers []byte
		if registers, err = client.WriteSingleCoil(5, 65280); err != nil {
			log.Fatalf("Error on write coils[5] = 1: %s", err)
		}
		log.Printf("Coils[5] = 1: %v", registers)
		if registers, err = client.WriteMultipleCoils(7, 3, []byte{1, 1, 0}); err != nil {
			log.Fatalf("Error on write coils[7:10] = {1, 1, 0}: %s", err)
		}
		log.Printf("Coils[7:10] = {1, 1, 0}: %v", registers)
		if registers, err = client.ReadCoils(5, 5); err != nil {
			log.Fatalf("Error on read coils[5:10]: %s", err)
		}
		log.Printf("Coils[5:10] = %v", registers)
		if registers, err = client.ReadDiscreteInputs(9, 11); err != nil {
			log.Fatalf("Error on read DI[9:20]: %s", err)
		}
		log.Printf("DI[9:20] = %v", registers)
		if registers, err = client.WriteSingleRegister(160, 84); err != nil {
			log.Fatalf("Error on write HR[160] = 84: %s", err)
		}
		log.Printf("HR[160] = 84: %v", registers)
		if registers, err = client.WriteMultipleRegisters(150, 7, []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59}); err != nil {
			log.Fatalf("Error on write HR[150:157] = {1, 18, 48, 53, 64, 57, 59}: %s", err)
		}
		log.Printf("HR[150:157] = {1, 18, 48, 53, 64, 57, 59}: %v", registers)
		if registers, err = client.ReadHoldingRegisters(150, 15); err != nil {
			log.Fatalf("Error on read HR[150:165]: %s", err)
		}
		log.Printf("HR[150:165] = %v", registers)
		if registers, err = client.ReadInputRegisters(4, 1); err != nil {
			log.Fatalf("Error on read IR[4]: %s", err)
		}
		log.Printf("IR[4] = %v", registers)
		if registers, err = client.ReadInputRegisters(4, 18); err != nil {
			log.Fatalf("Error on read IR[4:22]: %s", err)
		}
		log.Printf("IR[4:22] = %v", registers)
		if registers, err = client.ReadDiscreteInputs(11, 1); err != nil {
			log.Fatalf("Error on read DI[11]: %s", err)
		}
		log.Printf("DI[11] = %v", registers)
		if registers, err = client.ReadCoils(8, 1); err != nil {
			log.Fatalf("Error on read coils[8]: %s", err)
		}
		log.Printf("Coils[8] = %v", registers)
		if registers, err = client.ReadHoldingRegisters(153, 1); err != nil {
			log.Fatalf("Error on read HR[153]: %s", err)
		}
		log.Printf("HR[153] = %v", registers)
		handler.Close()
	}
}
