package trafficanalysis

import (
	"fmt"
	"modbus-emulator/conf"
	"modbus-emulator/src/traffic_analysis/structs"
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

func ParseDump() (history map[uint16]structs.ServerHistory, err error) {
	var currentHandle *pcap.Handle
	indexDictionary := make(map[structs.SlaveTransaction]int)
	history = make(map[uint16]structs.ServerHistory)
	for currentPhysicalPort, currentServerSocket := range conf.Ports {
		var currentHistory []structs.HistoryEvent
		var currentSlavesId []uint8
		for _, currentFilter := range []string{"dst", "src"} {
			if currentHandle, err = pcap.OpenOffline(fmt.Sprintf(`%s/%s/%s.pcapng`, conf.ModulePath, conf.DumpDirectoryPath, conf.WorkMode)); err != nil {
				err = fmt.Errorf("error on opening file: %s", err)
				return
			}
			if err = currentHandle.SetBPFFilter(fmt.Sprintf("tcp %s port %s", currentFilter, currentServerSocket.PortAddress)); err != nil {
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
				if conf.WorkMode == "rtu_over_tcp" {
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
				if !slices.Contains(currentSlavesId, currentHistoryEvent.Header.SlaveID) {
					currentSlavesId = append(currentSlavesId, currentHistoryEvent.Header.SlaveID)
				}
				if currentFilter == "dst" {
					currentHistoryEvent.Handshake = structs.Handshake{}
					currentHistoryEvent.Handshake.RequestUnmarshal(currentPayload)
					currentHistory = append(currentHistory, *currentHistoryEvent)
					indexDictionary[currentHistoryEvent.Header] = len(currentHistory) - 1
				} else {
					currentHistory[indexDictionary[currentHistoryEvent.Header]].Handshake.ResponseUnmarshal(currentPayload)
					currentHistory[indexDictionary[currentHistoryEvent.Header]].TransactionTime = currentPacket.Metadata().Timestamp
				}
			}
			currentHandle.Close()
		}
		history[currentPhysicalPort] = structs.ServerHistory{
			Transactions: currentHistory,
			Slaves:       currentSlavesId,
		}
	}
	return
}
