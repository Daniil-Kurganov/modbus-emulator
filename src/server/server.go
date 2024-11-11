package server

import (
	"fmt"
	"log"
	"time"

	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/utils"

	"github.com/tbrandon/mbserver"
)

var (
	history      []ta.History
	closeChannel = make(chan bool)
)

func emulate(server *mbserver.Server) {
	for currentIndex, currentHistoryEvent := range history {
		var timeEmulation time.Duration
		if currentIndex == len(history)-1 {
			timeEmulation = utils.FinishTime
		} else {
			timeEmulation = history[currentIndex+1].TransactionTime.Sub(currentHistoryEvent.TransactionTime)
		}
		currentHandshake := currentHistoryEvent.Handshake
		if !currentHandshake.Check() {
			log.Print("Current iteration check proccess was failed")
			continue
		}
		currentRequestData := currentHandshake.Request.MarshalData()
		var objectType, operation string
		switch currentHistoryEvent.Handshake.Request.GetHeader().FunctionType {
		case 5:
			server.Coils[currentRequestData.AddressStart[1]] = currentRequestData.Payload[0] + currentRequestData.Payload[1]
			objectType, operation = "coils", "simple writting"
			log.Print(server.Coils[:10])
		case 15:
			for currentIndex := int(currentRequestData.AddressStart[1]); currentIndex < int(currentRequestData.CheckField[1])+int(currentRequestData.AddressStart[1]); currentIndex++ {
				server.Coils[currentIndex] = currentRequestData.Payload[currentIndex-int(currentRequestData.AddressStart[1])]
			}
			objectType, operation = "coils", "multiple writting"
			log.Print(server.Coils[:10])
		case 6:
			server.HoldingRegisters[currentRequestData.AddressStart[1]] = uint16(currentRequestData.Payload[0]) + uint16(currentRequestData.Payload[1])
			objectType, operation = "HR", "simple writting"
			log.Print(server.HoldingRegisters[:10])
		case 16:
			counterIterations := 0
			for currentIndex := int(currentRequestData.AddressStart[1]); currentIndex < int(currentRequestData.CheckField[1])+int(currentRequestData.AddressStart[1]); currentIndex++ {
				server.HoldingRegisters[currentIndex] = uint16(currentRequestData.Payload[2*counterIterations]) + uint16(currentRequestData.Payload[2*counterIterations+1])
				counterIterations += 1
			}
			objectType, operation = "HR", "multiple writting"
			log.Print(server.HoldingRegisters[:10])
		}
		log.Printf("Current iteration:\n object type: %s\n operation: %s\n delay: %v\n\n", objectType, operation, timeEmulation)
		time.Sleep(timeEmulation)
	}
	log.Print("\nEnd of dump history file. Closing connection")
	closeChannel <- true
}

func Server() {
	var err error
	server := mbserver.NewServer()
	if err = server.ListenTCP(fmt.Sprintf("%s:%s", utils.ServerTCPHost, utils.ServerTCPPort)); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer server.Close()
	log.Printf("Start server on %s port", utils.ServerTCPPort)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
		}
	}()
	if history, err = ta.ParsePackets("workfiles", "HR", "write_32"); err != nil {
		log.Fatalf("Error on parsing dump history: %s", err)
	}
	go emulate(server)
	<-closeChannel
}
