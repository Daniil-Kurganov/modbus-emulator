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
		log.Printf("Current iteration: we will sleep %v", timeEmulation)
		time.Sleep(timeEmulation)
	}
	log.Print("End of dump history file. Closing connection")
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
