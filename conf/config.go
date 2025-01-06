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
	DumpSocketData struct {
		HostAddress string
		PortAddress string
		Protocol    string
	}
	DumpSocketsConfigData struct {
		DumpSocket string `toml:"DumpSocket" json:"dump_socket"`
		RealSocket string `toml:"RealSocket" json:"real_socket"`
		Protocol   string `toml:"Protocol" json:"protocol"`
	}
	TOMLConfig struct {
		ServerDefaultEmulateHost  string
		ServerHTTPServesocket     string
		ServerDefaultDumpPort     string
		FinishDelayTime           time.Duration
		DumpFilePath              string
		IsAutoParsingMode         bool
		EmulationPortAddressStart int
		OneTimeEmulation          bool
		DumpTimeLocation          string
		SimultaneouslyEmulation   bool
		DumpConfig                []DumpSocketsConfigData `toml:"DumpConfig"`
	}
)

var (
	Sockets                   map[string]DumpSocketData
	ServerDefaultEmulateHost  string
	ServerHTTPServesocket     string
	ServerDefaultDumpPort     string
	FinishDelayTime           time.Duration
	DumpFilePath              string
	IsAutoParsingMode         bool
	EmulationPortAddressStart uint16
	OneTimeEmulation          bool
	SimultaneouslyEmulation   bool
	DumpTimeLocation          *time.Location

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
	GenFileName   = "result_config.toml"
	GenFileTitles = struct {
		ServerDefaultEmulateHost  string
		ServerHTTPServesocket     string
		ServerDefaultDumpPort     string
		FinishDelayTime           string
		DumpFilePath              string
		IsAutoParsingMode         string
		EmulationPortAddressStart string
		OneTimeEmulation          string
		DumpTimeLocation          string
		SimultaneouslyEmulation   string
		DumpConfig                struct {
			Title string
			DumpSocketsConfigData
		}
	}{
		ServerDefaultEmulateHost:  "ServerDefaultEmulateHost",
		ServerHTTPServesocket:     "ServerHTTPServesocket",
		ServerDefaultDumpPort:     "ServerDefaultDumpPort",
		FinishDelayTime:           "FinishDelayTime",
		DumpFilePath:              "DumpFilePath",
		IsAutoParsingMode:         "IsAutoParsingMode",
		EmulationPortAddressStart: "EmulationPortAddressStart",
		OneTimeEmulation:          "OneTimeEmulation",
		DumpTimeLocation:          "DumpTimeLocation",
		SimultaneouslyEmulation:   "SimultaneouslyEmulation",
		DumpConfig: struct {
			Title string
			DumpSocketsConfigData
		}{
			Title: "[DumpConfig]",
			DumpSocketsConfigData: DumpSocketsConfigData{
				DumpSocket: " DumpSocket",
				RealSocket: " RealSocket",
				Protocol:   " Protocol",
			},
		},
	}
)

func init() {
	log.SetFlags(0)
	var err error
	var workDirectory string
	if workDirectory, err = os.Getwd(); err != nil {
		log.Fatalf("Error on configuration preprocessing: %s", err)
	}
	if strings.Contains(workDirectory, "tests") || strings.Contains(workDirectory, "utils") {
		return
	}
	var config TOMLConfig
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config.toml", workDirectory), &config); err != nil {
		log.Fatalf("Error on unmarshaling configuration: %s", err)
	}
	ServerDefaultEmulateHost = config.ServerDefaultEmulateHost
	ServerHTTPServesocket = config.ServerHTTPServesocket
	ServerDefaultDumpPort = config.ServerDefaultDumpPort
	FinishDelayTime = config.FinishDelayTime
	DumpFilePath = config.DumpFilePath
	IsAutoParsingMode = config.IsAutoParsingMode
	EmulationPortAddressStart = uint16(config.EmulationPortAddressStart)
	OneTimeEmulation = config.OneTimeEmulation
	if DumpTimeLocation, err = time.LoadLocation(config.DumpTimeLocation); err != nil {
		log.Fatalf("Error on parsing dump time location: %s", err)
	}
	SimultaneouslyEmulation = config.SimultaneouslyEmulation
	Sockets = make(map[string]DumpSocketData)
	if !IsAutoParsingMode {
		log.Print("Using manually work mode of parsing dump: using configuration list")
		for _, currentSocketData := range config.DumpConfig {
			var currentServePath string
			currentServerSocketData := DumpSocketData{Protocol: currentSocketData.Protocol}
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
	} else {
		log.Print("Using automatically work mode of parsing dump")
	}
}
