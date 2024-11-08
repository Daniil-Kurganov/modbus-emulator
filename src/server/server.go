package main

import (
	"log"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
	log.SetFlags(0)
	server := mbserver.NewServer()
	if err := server.ListenTCP("localhost:1502"); err != nil {
		log.Fatalf("Error on listening TCP: %s\n", err)
	}
	defer server.Close()
	log.Println("Start server on 1502 port")
	data := []uint16{8, 15, 39, 6}
	// server.InputRegisters[3] = 130
	// server.InputRegisters[5] = 101
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		server.InputRegisters[4] = data[counter]
		log.Printf("Set 4 IR to: %d", data[counter])
		counter++
		if counter >= len(data) {
			counter = 0
		}
	}
}
