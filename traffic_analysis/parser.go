package trafficanalysis

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"modbus-emulator/utils"
)

func ParsePackets(filename string) (err error) {
	var handle *pcap.Handle
	if handle, err = pcap.OpenOffline(fmt.Sprintf("./%s/%s.pcapng", utils.Foldername, filename)); err != nil {
		return fmt.Errorf("Error on opening file: %s\n", err)
	}
	defer handle.Close()
	if err = handle.SetBPFFilter("tcp src port 1502"); err != nil {
		return fmt.Errorf("Error on setting handle filter: %s\n", err)
	}
	packetsSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for currentPacket := range packetsSource.Packets() {
		currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
		currentPayload := currentTCPLayer.LayerPayload()
		if len(currentPayload) == 0 {
			continue
		}
		currentPacketNumber := currentPayload[0] + currentPayload[1]
		var currentProtocol string
		if currentPayload[2] == 0 && currentPayload[3] == 0 {
			currentProtocol = "modbus"
		} else {
			currentProtocol = "unknown"
		}
		currentBodyLength := currentPayload[4] + currentPayload[5]
		currentUnitID := currentPayload[6]
		var currentObjectType string
		switch currentPayload[7] {
		case 1:
			currentObjectType = "coils"
		case 2:
			currentObjectType = "discrete input"
		case 3:
			currentObjectType = "holding register"
		case 4:
			currentObjectType = "input register"
		default:
			currentObjectType = "unknown"
		}
		currentDataLength := currentPayload[8]
		currentData := currentPayload[9:]
		log.Printf("Transaction %v:\n protocol - %s;\n body length - %v;\n unit ID - %v;\n object - %s;\n data length - %v;\n data - %v\n\n",
			currentPacketNumber, currentProtocol, currentBodyLength, currentUnitID, currentObjectType, currentDataLength, currentData)
	}
	return
}
