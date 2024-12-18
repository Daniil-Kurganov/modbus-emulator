package main

import (
	"log"
	"modbus-emulator/conf"
	"modbus-emulator/src"
	ta "modbus-emulator/src/traffic_analysis"
	"modbus-emulator/src/traffic_analysis/structs"
	"sync"

	"golang.org/x/exp/maps"
)

func main() {
	log.SetFlags(0)
	var err error
	if conf.IsAutoParsingMode {
		if err = ta.SocketAutoAccumulation(); err != nil {
			log.Fatalf("Error on sockets auto accumulation: %s", err)
		}
		src.GenerateConfig()
	}
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	var waitGroup sync.WaitGroup
	for _, currentPhysicalSocket := range maps.Keys(conf.Sockets) {
		log.Print(currentPhysicalSocket)
		waitGroup.Add(1)
		go src.ServerInit(&waitGroup, currentPhysicalSocket, history[currentPhysicalSocket])
	}
	waitGroup.Wait()
	log.Print("All servers finished the work")
}
