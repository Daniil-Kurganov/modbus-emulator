package server

import (
	"log"
	"time"

	ta "modbus-emulator/src/traffic_analysis"

	"github.com/tbrandon/mbserver"
)

var (
	history []ta.History
	closeChannel = make(chan bool)
)

func emulate() () {
	for i := 0; i < 3; i++ {
		log.Println("Work")
		time.Sleep(time.Second)
	}
	closeChannel <- true
}

func Server() {
	server := mbserver.NewServer()
	if err := server.ListenTCP("localhost:1502"); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	log.Println("Start server on 1502 port")
	go func() {
		for {
			log.Println("Wait")
			time.Sleep(500 * time.Millisecond)
		}
	}()
	go emulate()
	a := <-closeChannel
	log.Println(a)
}
