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
		MarshalPayload() []uint16
		LogPrint()
	}
	Request interface {
		Packet
		MarshalAddress() []uint16
		MarshalQuantity() []uint16
	}
	Response interface {
		Packet
		GetFunctionID() uint16
	}

	HistoryEvent struct {
		TransactionID   string
		Handshake       Handshake
		TransactionTime time.Time
	}
	Handshake struct {
		Request  Request
		Response Response
	}
	EmulationData struct {
		FunctionID      uint16
		IsReadOperation bool
		Address         []uint16
		Quantity        []uint16
		Payload         []uint16
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
	if utils.WorkMode == "rtu_over_tcp" {
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
	if utils.WorkMode == "rtu_over_tcp" {
		functionID := payload[1]
		if slices.Contains([]byte{1, 2, 3, 4}, functionID) {
			hdhk.Response = new(RTUOverTCPReadResponse)
		} else if slices.Contains([]byte{5, 6}, functionID) {
			hdhk.Response = new(RTUOverTCPRequest123456Response56)
		} else if slices.Contains([]byte{15, 16}, functionID) {
			hdhk.Response = new(RTUOverTCPMultipleWriteResponse)
		} else {
			hdhk.Response = new(RTUOverTCPErrorResponse)
		}
	} else {
		hdhk.Response = new(TCPResponse)
	}
	hdhk.Response.Unmarshal(payload)
}

func (hdhk *Handshake) Marshal() (data EmulationData) {
	data.FunctionID = hdhk.Response.GetFunctionID()
	if slices.Contains([]uint16{5, 6, 15, 16}, data.FunctionID) {
		data.IsReadOperation = true
	}
	data.Address = hdhk.Request.MarshalAddress()
	data.Quantity = hdhk.Request.MarshalQuantity()
	if data.IsReadOperation {
		data.Payload = hdhk.Response.MarshalPayload()
	} else {
		data.Payload = hdhk.Request.MarshalPayload()
	}
	return
}

func (hdhk *Handshake) TransactionErrorCheck() bool {
	return slices.Contains([]uint16{1, 2, 3, 4, 5, 6, 15, 16}, hdhk.Response.GetFunctionID())
}
