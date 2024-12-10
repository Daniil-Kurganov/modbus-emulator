package main

import (
	"log"
	"modbus-emulator/conf"
	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"
	"time"

	"golang.org/x/exp/maps"
)

func main() {
	log.SetFlags(0)
	var history map[uint16]structs.ServerHistory
	var err error
	conf.WorkMode = "rtu_over_tcp"
	conf.DumpDirectoryPath = `pcapng_files/tests_files/multiple_ports`
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Printf("Ports: %v", maps.Keys(history))
	time.Sleep(time.Second)
	for currentPort, currentHistory := range history {
		var timeCount []time.Time
		log.Printf("Slaves for current port: %v", currentHistory.Slaves)
		for _, currentHistoryEvent := range currentHistory.Transactions {
			log.Printf("Current port: %d", currentPort)
			currentHistoryEvent.LogPrint()
			timeCount = append(timeCount, currentHistoryEvent.TransactionTime)
		}
		log.Printf("Port delay: %v", timeCount)
		time.Sleep(3 * time.Second)

	}
}
