package tcp

import (
	"fmt"
	"log"
	"time"

	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/utils"

	"github.com/tbrandon/mbserver"
)

var (
	Server       *mbserver.Server
	history      []ta.History
	closeChannel = make(chan bool)
)

func emulate() {
	for currentIndex, currentHistoryEvent := range history {
		var timeEmulation time.Duration
		if currentIndex == len(history)-1 {
			timeEmulation = utils.FinishTime
		} else {
			timeEmulation = history[currentIndex+1].TransactionTime.Sub(currentHistoryEvent.TransactionTime)
		}
		currentHistoryEvent.Print()
		currentHandshake := currentHistoryEvent.Handshake
		currentRequestData := currentHandshake.Request.MarshalData()
		var objectType, operation string
		switch currentHistoryEvent.Handshake.Request.GetHeader().FunctionType {
		case 1:
			currentPayload := currentHandshake.Response.MarshalData().Payload[0]
			currentAddress := currentRequestData.AddressStart[1] + currentRequestData.CheckField[1] - 1
			log.Print(currentAddress, currentPayload)
			if Server.Coils[currentAddress] != currentPayload {
				Server.Coils[currentAddress] = currentPayload
			}
			objectType, operation = "coils", "reading"
			log.Print(Server.Coils[:10])
		case 5:
			Server.Coils[currentRequestData.AddressStart[1]] = currentRequestData.Payload[0] + currentRequestData.Payload[1]
			objectType, operation = "coils", "simple writting"
			log.Print(Server.Coils[:10])
		case 15:
			for currentIndex := int(currentRequestData.AddressStart[1]); currentIndex < int(currentRequestData.CheckField[1])+int(currentRequestData.AddressStart[1]); currentIndex++ {
				Server.Coils[currentIndex] = currentRequestData.Payload[currentIndex-int(currentRequestData.AddressStart[1])]
			}
			objectType, operation = "coils", "multiple writting"
			log.Print(Server.Coils[:10])
		case 2:
			currentPayload := currentHandshake.Response.MarshalData().Payload[0]
			currentAddress := currentRequestData.AddressStart[1] + currentRequestData.CheckField[1] - 1
			log.Print(currentAddress, currentPayload)
			if Server.DiscreteInputs[currentAddress] != currentPayload {
				Server.DiscreteInputs[currentAddress] = currentPayload
			}
			objectType, operation = "DI", "reading"
			log.Print(Server.DiscreteInputs[:10])
		case 3:
			currentPayload := currentHandshake.Response.MarshalData().Payload
			currentFinishAddress := currentRequestData.AddressStart[1] + currentRequestData.CheckField[1]
			counterIterations := 0
			for currentAddress := currentRequestData.AddressStart[1]; currentAddress < currentFinishAddress; currentAddress++ {
				currentReadindBit := uint16(currentPayload[2*counterIterations]) + uint16(currentPayload[2*counterIterations+1])
				if Server.HoldingRegisters[currentAddress] != currentReadindBit {
					Server.HoldingRegisters[currentAddress] = currentReadindBit
				}
				counterIterations += 1
			}
			objectType, operation = "HR", "reading"
			log.Print(Server.HoldingRegisters[:10])
		case 6:
			Server.HoldingRegisters[currentRequestData.AddressStart[1]] = uint16(currentRequestData.Payload[0]) + uint16(currentRequestData.Payload[1])
			objectType, operation = "HR", "simple writting"
			log.Print(Server.HoldingRegisters[:10])
		case 16:
			counterIterations := 0
			for currentIndex := int(currentRequestData.AddressStart[1]); currentIndex < int(currentRequestData.CheckField[1])+int(currentRequestData.AddressStart[1]); currentIndex++ {
				Server.HoldingRegisters[currentIndex] = uint16(currentRequestData.Payload[2*counterIterations]) + uint16(currentRequestData.Payload[2*counterIterations+1])
				counterIterations += 1
			}
			objectType, operation = "HR", "multiple writting"
			log.Print(Server.HoldingRegisters[:10])
		case 4:
			currentPayload := currentHandshake.Response.MarshalData().Payload
			currentFinishAddress := currentRequestData.AddressStart[1] + currentRequestData.CheckField[1]
			counterIterations := 0
			for currentAddress := currentRequestData.AddressStart[1]; currentAddress < currentFinishAddress; currentAddress++ {
				currentReadindBit := uint16(currentPayload[2*counterIterations]) + uint16(currentPayload[2*counterIterations+1])
				if Server.InputRegisters[currentAddress] != currentReadindBit {
					Server.InputRegisters[currentAddress] = currentReadindBit
				}
				counterIterations += 1
			}
			objectType, operation = "IR", "reading"
			log.Print(Server.InputRegisters[:10])
		}
		log.Printf("Current iteration:\n object type: %s\n operation: %s\n delay: %v\n\n", objectType, operation, timeEmulation)
		time.Sleep(timeEmulation)
	}
	log.Print("\nEnd of dump history file. Closing connection")
	closeChannel <- true
}

func ServerInit() {
	var err error
	Server = mbserver.NewServer()
	if err = Server.ListenTCP(fmt.Sprintf("%s:%s", utils.ServerTCPHost, utils.ServerTCPPort)); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer Server.Close()
	log.Printf("Start server on %s port", utils.ServerTCPPort)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
		}
	}()
	if history, err = ta.ParsePackets(); err != nil {
		log.Fatalf("Error on parsing dump history: %s", err)
	}
	go emulate()
	<-closeChannel
}
