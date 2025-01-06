package src

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"os"
	"strings"

	tW "github.com/akiyosi/tomlwriter"
)

func GenerateConfig() (err error) {
	var newConfig []byte
	newConfig, _ = tW.WriteValue(fmt.Sprintf("\"%s\"", conf.ServerDefaultEmulateHost), newConfig, nil, conf.GenFileTitles.ServerDefaultEmulateHost, nil)
	newConfig, _ = tW.WriteValue(fmt.Sprintf("\"%s\"", conf.ServerHTTPServesocket), newConfig, nil, conf.GenFileTitles.ServerHTTPServesocket, nil)
	newConfig, _ = tW.WriteValue(fmt.Sprintf("\"%s\"", conf.ServerDefaultDumpPort), newConfig, nil, conf.GenFileTitles.ServerDefaultDumpPort, nil)
	newConfig, _ = tW.WriteValue(fmt.Sprintf("\"%s\"", conf.FinishDelayTime), newConfig, nil, conf.GenFileTitles.FinishDelayTime, nil)
	newConfig, _ = tW.WriteValue(fmt.Sprintf("'%s'", conf.DumpFilePath), newConfig, nil, conf.GenFileTitles.DumpFilePath, nil)
	newConfig, _ = tW.WriteValue(conf.IsAutoParsingMode, newConfig, nil, conf.GenFileTitles.IsAutoParsingMode, nil)
	newConfig, _ = tW.WriteValue(conf.EmulationPortAddressStart, newConfig, nil, conf.GenFileTitles.EmulationPortAddressStart, nil)
	newConfig, _ = tW.WriteValue(conf.OneTimeEmulation, newConfig, nil, conf.GenFileTitles.OneTimeEmulation, nil)
	newConfig, _ = tW.WriteValue(fmt.Sprintf("\"%s\"", conf.DumpTimeLocation), newConfig, nil, conf.GenFileTitles.DumpTimeLocation, nil)
	newConfig, _ = tW.WriteValue(conf.SimultaneouslyEmulation, newConfig, nil, conf.GenFileTitles.SimultaneouslyEmulation, nil)
	for currentEmulateSocket, currentDumpSocketData := range conf.Sockets {
		var currentDumpSocket, currentRealSocket string
		if currentDumpSocketData.PortAddress == conf.ServerDefaultDumpPort {
			currentDumpSocket = currentDumpSocketData.HostAddress
		} else {
			currentDumpSocket = fmt.Sprintf("%s:%s", currentDumpSocketData.HostAddress, currentDumpSocketData.PortAddress)
		}
		if currentEmulateSocket[:strings.Index(currentEmulateSocket, ":")] == conf.ServerDefaultEmulateHost {
			currentRealSocket = currentEmulateSocket[strings.Index(currentEmulateSocket, ":")+1:]
		} else {
			currentRealSocket = currentEmulateSocket
		}
		newConfig = append(newConfig, []byte(fmt.Sprintf("\n\n[%s]", conf.GenFileTitles.DumpConfig.Title))...)
		newConfig = append(newConfig, []byte(fmt.Sprintf("\n %s = \"%s\"", conf.GenFileTitles.DumpConfig.DumpSocket, currentDumpSocket))...)
		newConfig = append(newConfig, []byte(fmt.Sprintf("\n %s = \"%s\"", conf.GenFileTitles.DumpConfig.RealSocket, currentRealSocket))...)
		newConfig = append(newConfig, []byte(fmt.Sprintf("\n %s = \"%s\"", conf.GenFileTitles.DumpConfig.Protocol, currentDumpSocketData.Protocol))...)
	}
	var configFile *os.File
	if configFile, err = os.OpenFile(conf.GenFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		err = fmt.Errorf("error on creating new config file: %s", err)
		return
	}
	defer configFile.Close()
	if _, err = configFile.Write(newConfig); err != nil {
		err = fmt.Errorf("error on writing new config to file: %s", err)
		return
	}
	log.Printf("\nNew config file \"%s\" successfully written", conf.GenFileName)
	return
}
