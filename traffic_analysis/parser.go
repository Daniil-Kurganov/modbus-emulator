package trafficanalysis

import (
	"fmt"

	"modbus-emulator/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type TCPPacket struct {
	PacketNumber byte
	Protocol     string
	BodyLength   byte
	UnitID       byte
	ObjectType   byte
	DataLength   byte
	Data         []byte
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
		currentPacketResponse.Data = currentPayload[9:]
		packets = append(packets, currentPacketResponse)
	}
	return
}
