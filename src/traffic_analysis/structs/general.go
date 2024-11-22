package structs

import (
	"log"
	"modbus-emulator/src/utils"
	"slices"
	"time"
)

type (
	Packet interface {
		Unmarshal([]byte)
		LogPrint()
	}

	HistoryEvent struct {
		TransactionID   string
		Handshake       Handshake
		TransactionTime time.Time
	}
	Handshake struct {
		Request  Packet
		Response Packet
	}
)

func (hE *HistoryEvent) LogPrint() {
	log.Printf("\n\nTransaction № %v\n", hE.TransactionID)
	log.Println("\n Request:")
	hE.Handshake.Request.LogPrint()
	log.Println("\n Response:")
	hE.Handshake.Response.LogPrint()
	log.Printf("\n Transaction time: %v", hE.TransactionTime)
}

func (hdhk *Handshake) RequestUnmarshal(payload []byte) {
	if utils.Mode == "rtu_over_tcp" {
		functionID := payload[1]
		if slices.Contains([]byte{1, 2, 3, 4, 5, 6}, functionID) {
			hdhk.Request = new(RTUOverTCPRequest123456Response56)
		} else if slices.Contains([]byte{15, 16}, functionID) {
			hdhk.Request = new(RTUOverTCPMultipleWriteRequest)
		}
	} else {
		hdhk.Request = new(TCPRequest)
	}
	hdhk.Request.Unmarshal(payload)
}

func (hdhk *Handshake) ResponseUnmarshal(payload []byte) {
	if utils.Mode == "rtu_over_tcp" {
		functionID := payload[1]
		if slices.Contains([]byte{1, 2, 3, 4}, functionID) {
			hdhk.Response = new(RTUOverTCPReadResponse)
		} else if slices.Contains([]byte{5, 6}, functionID) {
			hdhk.Response = new(RTUOverTCPRequest123456Response56)
		} else if slices.Contains([]byte{15, 16}, functionID) {
			hdhk.Response = new(RTUOverTCPMultipleWriteResponse)
		}
	} else {
		hdhk.Response = new(TCPResponse)
	}
	hdhk.Response.Unmarshal(payload)
}
