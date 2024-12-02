package main

import (
	"log"
	"modbus-emulator/src"
	"sync"
)

func main() {
	log.SetFlags(0)
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go src.ServerInit(&waitGroup)
	go src.ServerInit(&waitGroup)
	waitGroup.Wait()
}
