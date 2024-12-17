package main

import (
	"io"
	"log"
	"os"
	"time"

	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"

	"golang.org/x/exp/maps"
)

func main() {
	log.SetFlags(0)
	var err error
	var logFile *os.File
	logFile, err = os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Printf("Sockets: %v", maps.Keys(history))
	time.Sleep(time.Second)
	for currentSocket, currentHistory := range history {
		log.Printf("Slaves for current port: %v", currentHistory.Slaves)
		for _, currentHistoryEvent := range currentHistory.Transactions {
			log.Printf("Current port: %s", currentSocket)
			currentHistoryEvent.LogPrint()
		}
		time.Sleep(3 * time.Second)
	}
}
