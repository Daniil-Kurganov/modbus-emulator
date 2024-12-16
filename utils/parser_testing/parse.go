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
	var history map[uint16]structs.ServerHistory
	// conf.WorkMode = "rtu_over_tcp"
	// conf.DumpDirectoryPath = `pcapng_files/tests_files/multiple_ports`
	// conf.Ports = map[uint16]conf.ServerSocket{
	// 	1502: {
	// 		HostAddress: "127.0.0.1",
	// 		PortAddress: "1502",
	// 	},
	// 	1503: {
	// 		HostAddress: "127.0.0.1",
	// 		PortAddress: "1503",
	// 	},
	// }
	if _, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Printf("Ports: %v", maps.Keys(history))
	time.Sleep(time.Second)
	for currentPort, currentHistory := range history {
		log.Printf("Slaves for current port: %v", currentHistory.Slaves)
		for _, currentHistoryEvent := range currentHistory.Transactions {
			log.Printf("Current port: %d", currentPort)
			currentHistoryEvent.LogPrint()
		}
		time.Sleep(3 * time.Second)
	}
}
