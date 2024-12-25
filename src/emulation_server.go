package src

import (
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"modbus-emulator/conf"
	"modbus-emulator/src/traffic_analysis/structs"

	mS "github.com/Daniil-Kurganov/modbus-server"
)

var (
	History            map[string]structs.ServerHistory
	IsEmulatingChannel chan (bool)
)

func ServerInit(waitGroup *sync.WaitGroup, servePath string) {
	var err error
	server := mS.NewServer()
	switch conf.Sockets[servePath].Protocol {
	case conf.Protocols.RTUOverTCP:
		if err = server.ListenRTUOverTCP(servePath); err != nil {
			log.Fatalf("Error on listening RTU over TCP: %s", err)
		}
	case conf.Protocols.TCP:
		if err = server.ListenTCP(servePath); err != nil {
			log.Fatalf("Error on listening TCP: %s", err)
		}
	default:
		log.Fatalf("Error: invalid servers's work mode: %s", conf.Sockets[servePath].Protocol)
	}
	log.Printf("Start server on %s, protocol: %s", servePath, conf.Sockets[servePath].Protocol)
	serverHistory := History[servePath]
	for _, currentSlaveId := range serverHistory.Slaves {
		server.InitSlave(currentSlaveId)
	}
	serverInfo := emulationServer{
		IsWorking: true,
		DumpSocketsConfigData: conf.DumpSocketsConfigData{
			DumpSocket: fmt.Sprintf("%s:%s", conf.Sockets[servePath].HostAddress, conf.Sockets[servePath].PortAddress),
			RealSocket: servePath,
			Protocol:   conf.Sockets[servePath].Protocol,
		},
		OneTimeEmulation: conf.OneTimeEmulation,
		StartTime:        serverHistory.Transactions[0].TransactionTime.String(),
		EndTime:          serverHistory.Transactions[len(serverHistory.Transactions)-1].TransactionTime.String(),
		CurrentTime:      "",
	}
	rewindChannel := make(chan int)
	emulationServers.readWriteMutex.Lock()
	emulationServers.serversData = append(emulationServers.serversData, serverInfo)
	emulationServers.servers = append(emulationServers.servers, server)
	emulationServers.rewindChannels = append(emulationServers.rewindChannels, rewindChannel)
	emulationServers.readWriteMutex.Unlock()
	emulationServers.readWriteMutex.RLock()
	serverID := len(emulationServers.serversData) - 1
	emulationServers.readWriteMutex.RUnlock()
	closeChannel := make(chan bool)
	go emulate(server, serverHistory.Transactions, closeChannel, serverID, rewindChannel)
	<-closeChannel
	close(closeChannel)
	server.Close()
	waitGroup.Done()
}

func emulate(server *mS.Server, history []structs.HistoryEvent, closeChannel chan (bool), serverID int, rewindChannel chan int) {
	if conf.SimultaneouslyEmulation {
		select {
		case <-server.ConnectionChanel:
			for counter := 0; counter < len(conf.Sockets)-1; counter++ {
				IsEmulatingChannel <- true
			}
		case <-IsEmulatingChannel:
			server.ConnectionChanel = nil
		}
	} else {
		log.Print("Waiting of client connection")
		<-server.ConnectionChanel
	}
	for {
		for currentIndex := 0; currentIndex < len(history); currentIndex++ {
			var currentHistoryEvent structs.HistoryEvent
			select {
			case transactionIndex := <-rewindChannel:
				log.Printf("Rewind (%d):\n %v", transactionIndex, history[transactionIndex])
				currentIndex = transactionIndex
			default:
			}
			currentHistoryEvent = history[currentIndex]
			emulationServers.readWriteMutex.Lock()
			emulationServers.serversData[serverID].CurrentTime = currentHistoryEvent.TransactionTime.String()
			emulationServers.readWriteMutex.Unlock()
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
		log.Print("\nEnd of dump history file.")
		emulationServers.readWriteMutex.Lock()
		defer emulationServers.readWriteMutex.Unlock()
		if emulationServers.serversData[serverID].OneTimeEmulation {
			log.Print("Emulation mode: one-time. Closing connection")
			emulationServers.serversData[serverID].IsWorking = false
			closeChannel <- true
			return
		}
		log.Print("Emulation mode: continuously. Starting new loop of emulation")
	}
}

func sliceUint16ToByte(source []uint16) (destination []byte) {
	for _, currentByte := range source {
		destination = append(destination, byte(currentByte))
	}
	return
}
