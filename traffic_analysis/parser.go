package trafficanalysis

import (
	"fmt"
	"log"

	"modbus-emulator/utils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func ParsePackets(filename string) (payload [][]byte, err error) {
	var handle *pcap.Handle
	if handle, err = pcap.OpenOffline(fmt.Sprintf("./%s/%s.pcapng", utils.Foldername, filename)); err != nil {
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
		payload = append(payload, currentPayload)
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
