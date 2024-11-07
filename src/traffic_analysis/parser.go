package trafficanalysis

import (
	"fmt"
	"log"
	"strconv"

	"modbus-emulator/src/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type History struct {
	TransactionID string
	Handshake     Handshake
}

func parsePacket(payload []byte, isRequest bool) (packet TCPPacket) {
	if isRequest {
		packet = new(TCPPacketRequest)
	} else {
		packet = new(TCPPacketResponse)
	}
	packet.UnmarshalHeader(payload)
	packet.UnmarshalData(payload)
	return
}

func transactionIDToKey(transcationID []byte) (key string) {
	for _, currentByte := range transcationID {
		key = fmt.Sprintf("%s-%s", key, strconv.Itoa(int(currentByte)))
	}
	key = key[1:]
	return
}

func ParsePackets(typeObject string, filename string) (history []History, err error) {
	var currentHandle *pcap.Handle
	indexDictionary := make(map[string]int)
	for _, currentFilter := range []string{"dst", "src"} {
		if currentHandle, err = pcap.OpenOffline(fmt.Sprintf("%s/%s/%s/%s.pcapng", utils.ModulePath, utils.Foldername, typeObject, filename)); err != nil {
			err = fmt.Errorf("error on opening file: %s", err)
			return
		}
		if err = currentHandle.SetBPFFilter(fmt.Sprintf("tcp %s port %s", currentFilter, utils.ServerTCPPort)); err != nil {
			err = fmt.Errorf("error on setting handle filter: %s", err)
			return
		}
		log.Print(currentFilter)
		currentPacketsSource := gopacket.NewPacketSource(currentHandle, currentHandle.LinkType())
		for currentPacket := range currentPacketsSource.Packets() {
			currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
			currentPayload := currentTCPLayer.LayerPayload()
			if len(currentPayload) == 0 {
				continue
			}
			log.Println(currentPayload)
			currentHistoryEvent := History{
				TransactionID: transactionIDToKey(currentPayload[:2]),
			}
			currentHandshake := Handshake{}
			if currentFilter == "dst" {
				currentHandshake.Request = parsePacket(currentPayload, true)
				currentHistoryEvent.Handshake = currentHandshake
				history = append(history, currentHistoryEvent)
				indexDictionary[currentHistoryEvent.TransactionID] = len(history) - 1
			} else {
				currentHandshake.Response = parsePacket(currentPayload, false)
				currentHistoryEvent.Handshake = currentHandshake
				history[indexDictionary[currentHistoryEvent.TransactionID]].Handshake.Response = currentHandshake.Response
			}
		}
		currentHandle.Close()
	}
	return
}
