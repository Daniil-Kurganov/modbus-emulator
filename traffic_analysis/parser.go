package trafficanalysis

import (
	"fmt"
	"log"
	"strconv"

	"modbus-emulator/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func parsePacket(payload []byte) (packet TCPPacket) {
	MBAPHeader := packet.UnmarshalHeader(payload)
	switch MBAPHeader.FunctionType {
	case 1:
		packet = new(TCPPacketRequest)
	}
	packet.UnmarshalHeader()
	packet = pkt
	return
}

func transactionIDToKey(transcationID []byte) (key string) {
	for _, currentByte := range transcationID {
		key = fmt.Sprintf("%s-%s", key, strconv.Itoa(int(currentByte)))
	}
	key = key[1:]
	return
}

func ParsePackets(typeObject string, filename string) (history map[string]Handshake, err error) {
	var currentHandle *pcap.Handle
	history = make(map[string]Handshake)
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
			currentTransactionID := transactionIDToKey(currentPayload[:2])
			currentHandshake, _ := history[currentTransactionID]
			// if currentFilter == "dst" {
			// 	currentRequest := packetParse(currentPayload, true)
			// 	currentRequest.TransactionID = currentPayload[:2]
			// 	currentHandshake.request = currentRequest
			// } else {
			// 	currentResponce := packetParse(currentPayload, false)
			// 	currentResponce.TransactionID = currentPayload[:2]
			// 	currentHandshake.responce = currentResponce
			// }
			history[currentTransactionID] = currentHandshake
		}
		currentHandle.Close()
	}
	return
}
