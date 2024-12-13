package main

import (
	"io"
	"log"
	"modbus-emulator/conf"
	"modbus-emulator/src"
	"os"
	"sync"

	"golang.org/x/exp/maps"
)

func main() {
	log.SetFlags(0)
	var logFile *os.File
	var err error
	logFile, err = os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	conf.WorkMode = "rtu_over_tcp"
	conf.DumpDirectoryPath = `pcapng_files/main_files/multiple_ports`
	conf.Ports = map[uint16]conf.ServerSocket{
		1502: {
			HostAddress: "127.0.0.1",
			PortAddress: "1502",
		},
		1503: {
			HostAddress: "127.0.0.1",
			PortAddress: "1503",
		},
	}
	var waitGroup sync.WaitGroup
	for _, currentPhysicalPort := range maps.Keys(conf.Ports) {
		log.Print(currentPhysicalPort)
		waitGroup.Add(1)
		go src.ServerInit(&waitGroup, currentPhysicalPort)
	}
	waitGroup.Wait()
	log.Print("All servers finished the work")
}
