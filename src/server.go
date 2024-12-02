package src

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"
	"modbus-emulator/src/utils"

	mS "github.com/Daniil-Kurganov/modbus-server"
)

func sliceUint16ToByte(source []uint16) (destination []byte) {
	for _, currentByte := range source {
		destination = append(destination, byte(currentByte))
	}
	return
}

func emulate(server *mS.Server, history []structs.HistoryEvent, closeChannel chan (bool)) {
	for currentIndex, currentHistoryEvent := range history {
		var timeEmulation time.Duration
		if currentIndex == len(history)-1 {
			timeEmulation = utils.FinishDelayTime
		} else {
			timeEmulation = history[currentIndex+1].TransactionTime.Sub(currentHistoryEvent.TransactionTime)
		}
		currentHistoryEvent.LogPrint()
		var currentObjectType, currentOperation string
		if currentHistoryEvent.Handshake.TransactionErrorCheck() {
			log.Print("Current transaction isn't valid")
			continue
		}
		var currentEmulationData structs.EmulationData
		var err error
		if currentEmulationData, err = currentHistoryEvent.Handshake.Marshal(); err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		currentRightBorder := int(currentEmulationData.Address + currentEmulationData.Quantity)
		switch currentEmulationData.FunctionID {
		case utils.Functions.CoilsRead:
			currentObjectType, currentOperation = "coils", "read"
			log.Printf("\n\n Before: Coils[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.Coils[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.Coils[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.Coils[currentIndex] = byte(currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)])
				}
			}
			log.Printf(" After: Coils[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.Coils[currentEmulationData.Address:currentRightBorder])
		case utils.Functions.DIRead:
			currentObjectType, currentOperation = "DI", "read"
			log.Printf("\n\n Before: DI[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.DiscreteInputs[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.DiscreteInputs[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.DiscreteInputs[currentIndex] = byte(currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)])
				}
			}
			log.Printf(" After: DI[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.DiscreteInputs[currentEmulationData.Address:currentRightBorder])
		case utils.Functions.HRRead:
			currentObjectType, currentOperation = "HR", "read"
			log.Printf("\n\n Before: HR[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.HoldingRegisters[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.HoldingRegisters[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.HoldingRegisters[currentIndex] = currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)]
				}
			}
			log.Printf(" After: HR[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.HoldingRegisters[currentEmulationData.Address:currentRightBorder])
		case utils.Functions.IRRead:
			currentObjectType, currentOperation = "IR", "read"
			log.Printf("\n\n Before: IR[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.InputRegisters[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.InputRegisters[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.InputRegisters[currentIndex] = currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)]
				}
			}
			log.Printf(" After: IR[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.InputRegisters[currentEmulationData.Address:currentRightBorder])
		case utils.Functions.CoilsSimpleWrite:
			currentObjectType, currentOperation = "coils", "simple write"
			log.Printf("\n\n Before: Coils[%d] = %d", currentEmulationData.Address, server.Coils[currentEmulationData.Address])
			server.Coils[currentEmulationData.Address] = byte(currentEmulationData.Payload[0])
			log.Printf(" After: Coils[%d] = %d", currentEmulationData.Address, server.Coils[currentEmulationData.Address])
		case utils.Functions.HRSimpleWrite:
			currentObjectType, currentOperation = "HR", "simple write"
			log.Printf("\n\n Before: HR[%d] = %d", currentEmulationData.Address, server.HoldingRegisters[currentEmulationData.Address])
			server.HoldingRegisters[currentEmulationData.Address] = currentEmulationData.Payload[0]
			log.Printf(" After: HR[%d] = %d", currentEmulationData.Address, server.HoldingRegisters[currentEmulationData.Address])
		case utils.Functions.CoilsMultipleWrite:
			currentObjectType, currentOperation = "coils", "multiple write"
			log.Printf("\n\n Before: Coils[%d:%d] = %v", currentEmulationData.Address, currentRightBorder, server.Coils[currentEmulationData.Address:currentRightBorder])
			for currentIndex := int(currentEmulationData.Address); currentIndex < int(currentRightBorder); currentIndex++ {
				server.Coils[currentIndex] = byte(currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)])
			}
			log.Printf(" After: Coils[%d:%d] = %v", currentEmulationData.Address, currentRightBorder, server.Coils[currentEmulationData.Address:currentRightBorder])
		case utils.Functions.HRMultipleWrite:
			currentObjectType, currentOperation = "HR", "multiple write"
			log.Printf("\n\n Before: HR[%d:%d] = %v", currentEmulationData.Address, currentRightBorder, server.HoldingRegisters[currentEmulationData.Address:currentRightBorder])
			for currentIndex := int(currentEmulationData.Address); currentIndex < int(currentRightBorder); currentIndex++ {
				server.HoldingRegisters[currentIndex] = currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)]
			}
			log.Printf(" After: HR[%d:%d] = %v", currentEmulationData.Address, currentRightBorder, server.HoldingRegisters[currentEmulationData.Address:currentRightBorder])
		}
		log.Printf("\nCurrent iteration:\n object type: %s\n operation: %s\n delay: %v\n\n", currentObjectType, currentOperation, timeEmulation)
		time.Sleep(timeEmulation)
	}
	log.Print("\nEnd of dump history file. Closing connection")
	closeChannel <- true
}

func ServerInit(waitGroup *sync.WaitGroup) {
	var err error
	server := mS.NewServer()
	servePath := fmt.Sprintf("%s:%s", utils.ServerTCPHost, utils.ServerTCPPort)
	if utils.WorkMode == "rtu_over_tcp" {
		if err = server.ListenRTUOverTCP(servePath); err != nil {
			log.Fatalf("Error on listening RTU over TCP: %s", err)
		}
	} else {
		if err = server.ListenTCP(servePath); err != nil {
			log.Fatalf("Error on listening TCP: %s", err)
		}
	}
	log.Printf("Start server on %s, workmode: %s", servePath, utils.WorkMode)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
		}
	}()
	var history []structs.HistoryEvent
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump history: %s", err)
	}
	closeChannel := make(chan bool)
	go emulate(server, history, closeChannel)
	<-closeChannel
	server.Close()
	waitGroup.Done()
}
