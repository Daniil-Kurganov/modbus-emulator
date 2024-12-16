package main

import (
	"log"
	"modbus-emulator/conf"
	"modbus-emulator/src"
	"sync"

	"golang.org/x/exp/maps"
)

func main() {
	log.SetFlags(0)
	var waitGroup sync.WaitGroup
	for _, currentPhysicalPort := range maps.Keys(conf.Ports) {
		log.Print(currentPhysicalPort)
		waitGroup.Add(1)
		go src.ServerInit(&waitGroup, currentPhysicalPort)
	}
	waitGroup.Wait()
	log.Print("All servers finished the work")
}
