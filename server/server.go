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
	data := []uint16{26, 99, 13, 40}
	counter := 0
	for {
		time.Sleep(2 * time.Second)
		server.InputRegisters[0] = data[counter]
		log.Printf("Set 0 input register to: %d", data[counter])
		counter++
		if counter >= len(data) {
			counter = 0
		}
	}
}
