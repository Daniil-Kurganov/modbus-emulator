package src

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"modbus-emulator/conf"
	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"

	mS "github.com/Daniil-Kurganov/modbus-server"
)

func sliceUint16ToByte(source []uint16) (destination []byte) {
	for _, currentByte := range source {
		destination = append(destination, byte(currentByte))
	}
	return
}

func emulate(server *mS.Server, history []structs.HistoryEvent, closeChannel chan (bool)) {
	log.Print("Waiting of client connection")
	<-server.ConnectionChanel
	for currentIndex, currentHistoryEvent := range history {
		var timeEmulation time.Duration
		if currentIndex == len(history)-1 {
			timeEmulation = conf.FinishDelayTime
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
		case conf.Functions.CoilsRead:
			currentObjectType, currentOperation = "coils", "read"
			log.Printf("\n\n Before: Coils[%d:%d] = %d",
				currentEmulationData.Address, currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentIndex] = byte(currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)])
				}
			}
			log.Printf(" After: Coils[%d:%d] = %d",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address:currentRightBorder])
		case conf.Functions.DIRead:
			currentObjectType, currentOperation = "DI", "read"
			log.Printf("\n\n Before: DI[%d:%d] = %d",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].DiscreteInputs[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.Slaves[currentHistoryEvent.Header.SlaveID].DiscreteInputs[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.Slaves[currentHistoryEvent.Header.SlaveID].DiscreteInputs[currentIndex] = byte(currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)])
				}
			}
			log.Printf(" After: DI[%d:%d] = %d",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].DiscreteInputs[currentEmulationData.Address:currentRightBorder])
		case conf.Functions.HRRead:
			currentObjectType, currentOperation = "HR", "read"
			log.Printf("\n\n Before: HR[%d:%d] = %d",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentIndex] = currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)]
				}
			}
			log.Printf(" After: HR[%d:%d] = %d",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address:currentRightBorder])
		case conf.Functions.IRRead:
			currentObjectType, currentOperation = "IR", "read"
			log.Printf("\n\n Before: IR[%d:%d] = %d", currentEmulationData.Address, currentRightBorder, server.Slaves[currentHistoryEvent.Header.SlaveID].InputRegisters[currentEmulationData.Address:currentRightBorder])
			if !reflect.DeepEqual(server.Slaves[currentHistoryEvent.Header.SlaveID].InputRegisters[currentEmulationData.Address:currentRightBorder], sliceUint16ToByte(currentEmulationData.Payload)) {
				for currentIndex := int(currentEmulationData.Address); currentIndex < currentRightBorder; currentIndex++ {
					server.Slaves[currentHistoryEvent.Header.SlaveID].InputRegisters[currentIndex] = currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)]
				}
			}
			log.Printf(" After: IR[%d:%d] = %d",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].InputRegisters[currentEmulationData.Address:currentRightBorder])
		case conf.Functions.CoilsSimpleWrite:
			currentObjectType, currentOperation = "coils", "simple write"
			log.Printf("\n\n Before: Coils[%d] = %d",
				currentEmulationData.Address,
				server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address])
			server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address] = byte(currentEmulationData.Payload[0])
			log.Printf(" After: Coils[%d] = %d", currentEmulationData.Address, server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address])
		case conf.Functions.HRSimpleWrite:
			currentObjectType, currentOperation = "HR", "simple write"
			log.Printf("\n\n Before: HR[%d] = %d", currentEmulationData.Address, server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address])
			server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address] = currentEmulationData.Payload[0]
			log.Printf(" After: HR[%d] = %d", currentEmulationData.Address, server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address])
		case conf.Functions.CoilsMultipleWrite:
			currentObjectType, currentOperation = "coils", "multiple write"
			log.Printf("\n\n Before: Coils[%d:%d] = %v",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address:currentRightBorder])
			for currentIndex := int(currentEmulationData.Address); currentIndex < int(currentRightBorder); currentIndex++ {
				server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentIndex] = byte(currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)])
			}
			log.Printf(" After: Coils[%d:%d] = %v",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].Coils[currentEmulationData.Address:currentRightBorder])
		case conf.Functions.HRMultipleWrite:
			currentObjectType, currentOperation = "HR", "multiple write"
			log.Printf("\n\n Before: HR[%d:%d] = %v",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address:currentRightBorder])
			for currentIndex := int(currentEmulationData.Address); currentIndex < int(currentRightBorder); currentIndex++ {
				server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentIndex] = currentEmulationData.Payload[currentIndex-int(currentEmulationData.Address)]
			}
			log.Printf(" After: HR[%d:%d] = %v",
				currentEmulationData.Address,
				currentRightBorder,
				server.Slaves[currentHistoryEvent.Header.SlaveID].HoldingRegisters[currentEmulationData.Address:currentRightBorder])
		}
		log.Printf("\nCurrent iteration:\n slave ID: %d\n object type: %s\n operation: %s\n delay: %v\n\n",
			currentHistoryEvent.Header.SlaveID,
			currentObjectType,
			currentOperation,
			timeEmulation)
		time.Sleep(timeEmulation)
	}
	log.Print("\nEnd of dump history file. Closing connection")
	closeChannel <- true
}

func ServerInit(waitGroup *sync.WaitGroup, physicalPort string) {
	var err error
	server := mS.NewServer()
	servePath := fmt.Sprintf("%s:%s", conf.ServerTCPHost, physicalPort)
	switch conf.Ports[physicalPort].WorkMode {
	case "rtu_over_tcp":
		if err = server.ListenRTUOverTCP(servePath); err != nil {
			log.Fatalf("Error on listening RTU over TCP: %s", err)
		}
	case "tcp":
		if err = server.ListenTCP(servePath); err != nil {
			log.Fatalf("Error on listening TCP: %s", err)
		}
	default:
		log.Fatalf("Error: invalid servers's work mode: %s", conf.Ports[physicalPort].WorkMode)
	}
	log.Printf("Start server on %s, work mode: %s", servePath, conf.Ports[physicalPort].WorkMode)
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump history: %s", err)
	}
	for _, currentSlaveId := range history[physicalPort].Slaves {
		server.InitSlave(currentSlaveId)
	}
	closeChannel := make(chan bool)
	go emulate(server, history[physicalPort].Transactions, closeChannel)
	<-closeChannel
	close(closeChannel)
	server.Close()
	waitGroup.Done()
}
