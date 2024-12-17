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
	for _, currentPhysicalSocket := range maps.Keys(conf.Sockets) {
		log.Print(currentPhysicalSocket)
		waitGroup.Add(1)
		go src.ServerInit(&waitGroup, currentPhysicalSocket)
	}
	waitGroup.Wait()
	log.Print("All servers finished the work")
}
