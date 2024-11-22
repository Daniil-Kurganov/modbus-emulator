package trafficanalysis

import (
	"fmt"
	"modbus-emulator/src/traffic_analysis/structs"
	"modbus-emulator/src/utils"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func TCPTransactionIDParsing(transcationID []byte) (key string) {
	for _, currentByte := range transcationID {
		key = fmt.Sprintf("%s-%s", key, strconv.Itoa(int(currentByte)))
	}
	key = key[1:]
	return
}

func ParseDump() (history []structs.HistoryEvent, err error) {
	var currentHandle *pcap.Handle
	indexDictionary := make(map[string]int)
	for _, currentFilter := range []string{"dst", "src"} {
		if currentHandle, err = pcap.OpenOffline(fmt.Sprintf(`%s/%s/%s.pcapng`, utils.ModulePath, utils.DumpDirectoryPath, utils.Mode)); err != nil {
			err = fmt.Errorf("error on opening file: %s", err)
			return
		}
		if err = currentHandle.SetBPFFilter(fmt.Sprintf("tcp %s port %s", currentFilter, utils.ServerTCPPort)); err != nil {
			err = fmt.Errorf("error on setting handle filter: %s", err)
			return
		}
		currentPacketsSource := gopacket.NewPacketSource(currentHandle, currentHandle.LinkType())
		counterTransaction := 1
		for currentPacket := range currentPacketsSource.Packets() {
			currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
			currentPayload := currentTCPLayer.LayerPayload()
			if len(currentPayload) == 0 {
				continue
			}
			currentHistoryEvent := new(structs.HistoryEvent)
			if utils.Mode == "rtu_over_tcp" {
				currentHistoryEvent.TransactionID = strconv.Itoa(counterTransaction)
				counterTransaction += 1
			} else {
				currentHistoryEvent.TransactionID = TCPTransactionIDParsing(currentPayload[:2])
			}
			if currentFilter == "dst" {
				currentHistoryEvent.Handshake = structs.Handshake{}
				currentHistoryEvent.Handshake.RequestUnmarshal(currentPayload)
				history = append(history, *currentHistoryEvent)
				indexDictionary[currentHistoryEvent.TransactionID] = len(history) - 1
			} else {
				history[indexDictionary[currentHistoryEvent.TransactionID]].Handshake.ResponseUnmarshal(currentPayload)
				history[indexDictionary[currentHistoryEvent.TransactionID]].TransactionTime = currentPacket.Metadata().Timestamp
			}
		}
		currentHandle.Close()
	}
	return
}
