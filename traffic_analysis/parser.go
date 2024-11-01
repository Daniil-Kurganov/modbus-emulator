package trafficanalysis

import (
	"fmt"

	"modbus-emulator/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type TCPPacket struct {
	packetNumber byte
	protocol     string
	bodyLength   byte
	unitID       byte
	objectType   byte
	dataLength   byte
	data         []byte
}

func ParsePackets(filename string) (packets []TCPPacket, err error) {
	var handle *pcap.Handle
	if handle, err = pcap.OpenOffline(fmt.Sprintf("%s/%s/%s.pcapng", utils.ModulePath, utils.Foldername, filename)); err != nil {
		err = fmt.Errorf("error on opening file: %s", err)
		return
	}
	defer handle.Close()
	if err = handle.SetBPFFilter("tcp src port 1502"); err != nil {
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
		var currentPacketResponse TCPPacket
		currentPacketResponse.packetNumber = currentPayload[0] + currentPayload[1]
		if currentPayload[2] == 0 && currentPayload[3] == 0 {
			currentPacketResponse.protocol = "modbus"
		} else {
			currentPacketResponse.protocol = "unknown"
		}
		currentPacketResponse.bodyLength = currentPayload[4] + currentPayload[5]
		currentPacketResponse.unitID = currentPayload[6]
		currentPacketResponse.objectType = currentPayload[7]
		currentPacketResponse.dataLength = currentPayload[8]
		currentPacketResponse.data = currentPayload[9:]
		packets = append(packets, currentPacketResponse)
	}
	return
}
