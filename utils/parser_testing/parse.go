package main

import (
	"log"
	"time"

	// "modbus-emulator/conf"
	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"

	"golang.org/x/exp/maps"
)

func main() {
	log.SetFlags(0)
	var history map[uint16]structs.ServerHistory
	var err error
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
	if history, err = ta.ParseDump(); err != nil {
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
		return
	}
}
