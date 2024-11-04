package trafficanalysis

import (
	"fmt"
	"slices"

	"modbus-emulator/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func packetParse(payload []byte, isRequset bool) (packet TCPPacket) {
	if payload[2] == 0 && payload[3] == 0 {
		packet.Protocol = "modbus"
	} else {
		packet.Protocol = "unknown"
	}
	packet.BodyLength = payload[4] + payload[5]
	packet.UnitID = payload[6]
	packet.FunctionType = payload[7]
	packet.AddressStart = payload[8:10]
	if slices.Contains([]byte{1, 2, 3, 4}, packet.FunctionType) {
		if isRequset {
			packet.DataPayload = &readRequest{
				operationCode: 0,
				payload:       payload[10:],
			}
		} else {
			var payloadData [][]byte
			var workData []byte
			for currentIndex, currentBit := range payload[11:] {
				if currentIndex%2 == 0 {
					workData = []byte{currentBit}
				} else {
					workData = append(workData, currentBit)
					payloadData = append(payloadData, workData)
				}
			}
			packet.DataPayload = &readResponce{
				operationCode: 0,
				payload: struct {
					numberBits byte
					data       [][]byte
				}{payload[10], payloadData},
			}
		}
	}
	return
}

func ParsePackets(typeObject string, filename string) (packets []TCPPacket, err error) {
	var handle *pcap.Handle
	if handle, err = pcap.OpenOffline(fmt.Sprintf("%s/%s/%s/%s.pcapng", utils.ModulePath, utils.Foldername, typeObject, filename)); err != nil {
		err = fmt.Errorf("error on opening file: %s", err)
		return
	}
	defer handle.Close()
	history := make(map[byte]Handshake)
	for _, currentFilter := range []string{"dst", "src"} {
		if err = handle.SetBPFFilter(fmt.Sprintf("tcp %s port %s", currentFilter, utils.ServerTCPPort)); err != nil {
			err = fmt.Errorf("error on setting handle filter: %s", err)
			return
		}
		currentPacketsSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for currentPacket := range currentPacketsSource.Packets() {
			currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
			currentPayload := currentTCPLayer.LayerPayload()
			if len(currentPayload) == 0 {
				continue
			}
			currentPacketNumber := currentPayload[0] + currentPayload[1]
			currentHandshake, _ := history[currentPacketNumber]
			if currentFilter == "dst" {
				currentRequest := packetParse(currentPayload, true)
				currentRequest.PacketNumber = currentPacketNumber
				currentHandshake
			} else {

			}
		}
	}

	// for currentPacket := range packetsSource.Packets() {
	// 	currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
	// 	currentPayload := currentTCPLayer.LayerPayload()
	// 	if len(currentPayload) == 0 {
	// 		continue
	// 	}
	// 	log.Println(currentPayload)
	// 	var currentPacketResponse TCPPacket
	// 	currentPacketResponse.PacketNumber = currentPayload[0] + currentPayload[1]
	//
	// 	packets = append(packets, packetResponse)
	// }
	return
}
