package trafficanalysis

import (
	"fmt"
	"modbus-emulator/src/traffic_analysis/structs"
	"modbus-emulator/src/utils"
	"slices"
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

func ParseDump() (history []structs.HistoryEvent, slavesId []uint8, err error) {
	var currentHandle *pcap.Handle
	indexDictionary := make(map[structs.SlaveTransaction]int)
	for _, currentFilter := range []string{"dst", "src"} {
		if currentHandle, err = pcap.OpenOffline(fmt.Sprintf(`%s/%s/%s.pcapng`, utils.ModulePath, utils.DumpDirectoryPath, utils.WorkMode)); err != nil {
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
			if utils.WorkMode == "rtu_over_tcp" {
				currentHistoryEvent.Header = structs.SlaveTransaction{
					SlaveID:       uint8(currentPayload[0]),
					TransactionID: strconv.Itoa(counterTransaction),
				}
				counterTransaction += 1
			} else {
				currentHistoryEvent.Header = structs.SlaveTransaction{
					SlaveID:       uint8(currentPayload[6]),
					TransactionID: TCPTransactionIDParsing(currentPayload[:2]),
				}
			}
			if !slices.Contains(slavesId, currentHistoryEvent.Header.SlaveID) {
				slavesId = append(slavesId, currentHistoryEvent.Header.SlaveID)
			}
			if currentFilter == "dst" {
				currentHistoryEvent.Handshake = structs.Handshake{}
				currentHistoryEvent.Handshake.RequestUnmarshal(currentPayload)
				history = append(history, *currentHistoryEvent)
				indexDictionary[currentHistoryEvent.Header] = len(history) - 1
			} else {
				history[indexDictionary[currentHistoryEvent.Header]].Handshake.ResponseUnmarshal(currentPayload)
				history[indexDictionary[currentHistoryEvent.Header]].TransactionTime = currentPacket.Metadata().Timestamp
			}
		}
		currentHandle.Close()
	}
	return
}
