package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/simonvetter/modbus"
)

func client(wG *sync.WaitGroup, port uint16) {
	var err error
	var client *modbus.ModbusClient
	if client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("rtuovertcp://localhost:%d", port),
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
	// var sCoils, sDI bool
	var mHR, mIR []uint16
	// var sHR, sIR uint16``
	for _, currentUnitID := range []uint8{1, 2, 3} {
		if err = client.SetUnitId(currentUnitID); err != nil {
			log.Fatalf("Error on setting uint ID")
		}
		log.Printf("\nSet uint ID = %d", currentUnitID)
		// if err = client.WriteCoil(5, true); err != nil {
		// 	log.Fatalf("Error on writting coils[5] = 1: %s", err)
		// }
		if err = client.WriteCoils(7, []bool{true, true, false}); err != nil {
			log.Fatalf("Error on writting coils[7:10] = {1, 1, 0}: %s", err)
		}

		if coils, err = client.ReadCoils(5, 5); err != nil {
			log.Fatalf("Error on reading coils[5:10]: %s", err)
		}
		log.Printf("Coils[5:10] = %v", coils)
		if DI, err = client.ReadDiscreteInputs(9, 11); err != nil {
			log.Fatalf("Error on reading DI[9:20]: %s", err)
		}
		log.Printf("DI[9:20] = %v", DI)
		// if err = client.WriteRegister(160, 84); err != nil {
		// 	log.Fatalf("Error on write HR[160] = 84")
		// }
		if err = client.WriteRegisters(150, []uint16{1, 18, 48, 53, 64, 57, 59}); err != nil {
			log.Fatalf("Error on write HR[150:157] = {1, 18, 48, 53, 64, 57, 59}: %s", err)
		}
		if mHR, err = client.ReadRegisters(150, 15, modbus.HOLDING_REGISTER); err != nil {
			log.Fatalf("Error on read HR[150:165]: %s", err)
		}
		log.Printf("HR[150:165] = %v", mHR)
		// if sIR, err = client.ReadRegister(4, modbus.INPUT_REGISTER); err != nil {
		// 	log.Fatalf("Error on read IR[4]: %s", err)
		// }
		// log.Printf("IR[4] = %d", sIR)
		if mIR, err = client.ReadRegisters(4, 18, modbus.INPUT_REGISTER); err != nil {
			log.Fatalf("Error on read IR[4:22]: %s", err)
		}
		log.Printf("IR[4:22] = %v", mIR)
		// if sDI, err = client.ReadDiscreteInput(11); err != nil {
		// 	log.Fatalf("Error on read DI[11]: %s", err)
		// }
		// log.Printf("DI[11] = %t", sDI)
		// if sCoils, err = client.ReadCoil(8); err != nil {
		// 	log.Fatalf("Error on read coils[8]: %s", err)
		// }
		// log.Printf("coils[8] = %t", sCoils)
		// if sHR, err = client.ReadRegister(153, modbus.HOLDING_REGISTER); err != nil {
		// 	log.Fatalf("Error on read HR[153]: %s", err)
		// }
		// log.Printf("HR[153] = %d", sHR)
	}
	wG.Done()
}

func main() {
	log.SetFlags(0)
	var wG sync.WaitGroup
	wG.Add(2)
	go client(&wG, 1502)
	go client(&wG, 1503)
	wG.Wait()
}
