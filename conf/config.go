package conf

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type (
	ServerSocketData struct {
		HostAddress string
		PortAddress string
		Protocol    string
	}
	DumpSocketsConfigData struct {
		DumpSocket string `toml:"DumpSocket"`
		RealSocket string `toml:"RealSocket"`
		Protocol   string `toml:"Protocol"`
	}
	TOMLConfig struct {
		DumpConfig               []DumpSocketsConfigData `toml:"DumpConfig"`
		ServerDefaultEmulateHost string
		ServerDefaultDumpPort    string
		FinishDelayTime          time.Duration
		DumpFilePath             string
	}
)

var (
	Sockets                  map[string]ServerSocketData
	ServerDefaultEmulateHost string
	ServerDefaultDumpPort    string
	FinishDelayTime          time.Duration
	DumpFilePath             string

	Functions = struct {
		CoilsRead          uint16
		DIRead             uint16
		HRRead             uint16
		IRRead             uint16
		CoilsSimpleWrite   uint16
		HRSimpleWrite      uint16
		CoilsMultipleWrite uint16
		HRMultipleWrite    uint16
	}{
		CoilsRead:          1,
		DIRead:             2,
		HRRead:             3,
		IRRead:             4,
		CoilsSimpleWrite:   5,
		HRSimpleWrite:      6,
		CoilsMultipleWrite: 15,
		HRMultipleWrite:    16,
	}
	Protocols = struct {
		RTUOverTCP string
		TCP        string
	}{
		RTUOverTCP: "rtu_over_tcp",
		TCP:        "tcp",
	}
)

func init() {
	log.SetFlags(0)
	var err error
	var workDirectory string
	if workDirectory, err = os.Getwd(); err != nil {
		log.Fatalf("Error on configuration preprocessing: %s", err)
	}
	if strings.Contains(workDirectory, "tests") {
		return
	}
	var config TOMLConfig
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config.toml", workDirectory), &config); err != nil {
		log.Fatalf("Error on unmarshaling configuration: %s", err)
	}
	ServerDefaultEmulateHost = config.ServerDefaultEmulateHost
	ServerDefaultDumpPort = config.ServerDefaultDumpPort
	FinishDelayTime = config.FinishDelayTime
	DumpFilePath = config.DumpFilePath
	Sockets = make(map[string]ServerSocketData)
	for _, currentSocketData := range config.DumpConfig {
		var currentServePath string
		currentServerSocketData := ServerSocketData{Protocol: currentSocketData.Protocol}
		if currentSepIndex := strings.Index(currentSocketData.DumpSocket, ":"); currentSepIndex == -1 {
			currentServerSocketData.HostAddress = currentSocketData.DumpSocket
			currentServerSocketData.PortAddress = ServerDefaultDumpPort
		} else {
			currentServerSocketData.HostAddress = currentSocketData.DumpSocket[:currentSepIndex]
			currentServerSocketData.PortAddress = currentSocketData.DumpSocket[currentSepIndex+1:]
		}
		if currentSepIndex := strings.Index(currentSocketData.RealSocket, ":"); currentSepIndex == -1 {
			currentServePath = fmt.Sprintf("%s:%s", ServerDefaultEmulateHost, currentSocketData.RealSocket)
		} else {
			currentServePath = currentSocketData.RealSocket
		}
		Sockets[currentServePath] = currentServerSocketData
	}
}
