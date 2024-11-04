package trafficanalysis

import (
	"fmt"
	"log"

	"modbus-emulator/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func ParsePackets(typeObject string, filename string, r string) (packets []TCPPacket, err error) {
	var handle *pcap.Handle
	if handle, err = pcap.OpenOffline(fmt.Sprintf("%s/%s/%s/%s.pcapng", utils.ModulePath, utils.Foldername, typeObject, filename)); err != nil {
		err = fmt.Errorf("error on opening file: %s", err)
		return
	}
	defer handle.Close()
	if err = handle.SetBPFFilter(fmt.Sprintf("tcp %s port 1502", r)); err != nil {
		err = fmt.Errorf("error on setting handle filter: %s", err)
		return
	}
	packetsSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for currentPacket := range packetsSource.Packets() {
		currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
		currentPayload := currentTCPLayer.LayerPayload()
		if len(currentPayload) == 0 {
			continue
		}
		log.Println(currentPayload)
		var currentPacketResponse TCPPacket
		currentPacketResponse.PacketNumber = currentPayload[0] + currentPayload[1]
		if currentPayload[2] == 0 && currentPayload[3] == 0 {
			currentPacketResponse.Protocol = "modbus"
		} else {
			currentPacketResponse.Protocol = "unknown"
		}
		currentPacketResponse.BodyLength = currentPayload[4] + currentPayload[5]
		currentPacketResponse.UnitID = currentPayload[6]
		currentPacketResponse.ObjectType = currentPayload[7]
		currentPacketResponse.DataLength = currentPayload[8]
		// currentPacketResponse.DataPayload = currentPayload[9:]
		packets = append(packets, currentPacketResponse)
	}
	return
}
