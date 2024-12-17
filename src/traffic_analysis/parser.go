package trafficanalysis

import (
	"fmt"
	"log"
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

func ParseDump() (history map[string]structs.ServerHistory, err error) {
	var currentHandle *pcap.Handle
	history = make(map[string]structs.ServerHistory)
	for currentPhysicalSocket, currentServerSocketData := range conf.Sockets {
		if currentHandle, err = pcap.OpenOffline(fmt.Sprintf(`%s.pcapng`, conf.DumpFilePath)); err != nil {
			if currentHandle, err = pcap.OpenOffline(fmt.Sprintf(`%s.pcap`, conf.DumpFilePath)); err != nil {
				err = fmt.Errorf("error on opening file: %s", err)
				return
			}

		}
		if err = currentHandle.SetBPFFilter(fmt.Sprintf("host %s and tcp port %s",
			currentServerSocketData.HostAddress, currentServerSocketData.PortAddress)); err != nil {
			err = fmt.Errorf("error on setting handle filter: %s", err)
			return
		}
		currentPacketsSource := gopacket.NewPacketSource(currentHandle, currentHandle.LinkType())
		var currentHistory []structs.HistoryEvent
		var currentSlavesId []uint8
		var rtuOverTCPTransactionDictionary map[uint8]int
		if currentServerSocketData.Protocol == conf.Protocols.RTUOverTCP {
			rtuOverTCPTransactionDictionary = make(map[uint8]int)
		}
		for currentPacket := range currentPacketsSource.Packets() {
			currentTCPLayer := currentPacket.Layer(layers.LayerTypeTCP)
			currentPayload := currentTCPLayer.LayerPayload()
			if len(currentPayload) == 0 {
				continue
			}
			currentPacketIsRequest := currentPacket.TransportLayer().TransportFlow().Dst().String() == currentServerSocketData.PortAddress
			if !currentPacketIsRequest {
				if len(currentHistory) == 0 {
					continue
				}
				if currentHistory[len(currentHistory)-1].Handshake.Response != nil {
					if currentServerSocketData.Protocol == conf.Protocols.RTUOverTCP {
						rtuOverTCPTransactionDictionary[currentHistory[len(currentHistory)-1].Header.SlaveID] -= 1
					}
					currentHistory = currentHistory[:len(currentHistory)-1]
					continue
				}
				currentHistory[len(currentHistory)-1].Handshake.ResponseUnmarshal(currentServerSocketData.Protocol, currentPayload)
				currentHistory[len(currentHistory)-1].TransactionTime = currentPacket.Metadata().Timestamp
			} else {
				if len(currentHistory) != 0 && currentHistory[len(currentHistory)-1].Handshake.Response == nil {
					if currentServerSocketData.Protocol == conf.Protocols.RTUOverTCP {
						rtuOverTCPTransactionDictionary[currentHistory[len(currentHistory)-1].Header.SlaveID] -= 1
					}
					currentHistory = currentHistory[:len(currentHistory)-1]
					continue
				}
				currentHistoryEvent := new(structs.HistoryEvent)
				switch currentServerSocketData.Protocol {
				case conf.Protocols.RTUOverTCP:
					currentSlaveId := uint8(currentPayload[0])
					if _, ok := rtuOverTCPTransactionDictionary[currentSlaveId]; !ok {
						rtuOverTCPTransactionDictionary[currentSlaveId] = 1
					} else {
						rtuOverTCPTransactionDictionary[currentSlaveId] += 1
					}
					currentHistoryEvent.Header = structs.SlaveTransaction{
						SlaveID:       currentSlaveId,
						TransactionID: strconv.Itoa(rtuOverTCPTransactionDictionary[currentSlaveId]),
					}
				case conf.Protocols.TCP:
					currentHistoryEvent.Header = structs.SlaveTransaction{
						SlaveID:       uint8(currentPayload[6]),
						TransactionID: TCPTransactionIDParsing(currentPayload[:2]),
					}
				default:
					log.Fatalf("Error on parsing dump: %+v has invalid protocol", currentServerSocketData)
				}
				if !slices.Contains(currentSlavesId, currentHistoryEvent.Header.SlaveID) {
					currentSlavesId = append(currentSlavesId, currentHistoryEvent.Header.SlaveID)
				}
				currentHistoryEvent.Handshake.RequestUnmarshal(currentServerSocketData.Protocol, currentPayload)
				currentHistory = append(currentHistory, *currentHistoryEvent)
			}
		}
		currentHandle.Close()
		currentPortHistory := structs.ServerHistory{
			Transactions: currentHistory,
			Slaves:       currentSlavesId,
		}
		currentPortHistory.SelfClean()
		history[currentPhysicalSocket] = currentPortHistory
	}
	return
}
