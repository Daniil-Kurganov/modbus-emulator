package structs

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"slices"
	"strconv"
	"strings"
	"time"
)

type (
	Packet interface {
		Unmarshal([]byte)
		MarshalPayload() ([]uint16, error)
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

	SlaveTransaction struct {
		SlaveID       uint8
		TransactionID string
	}
	HistoryEvent struct {
		Header          SlaveTransaction
		Handshake       Handshake
		TransactionTime time.Time
	}
	ServerHistory struct {
		Transactions []HistoryEvent
		Slaves       []uint8
	}
	Handshake struct {
		Request  Request
		Response Response
	}
	EmulationData struct {
		FunctionID      uint16
		IsReadOperation bool
		Address         uint16
		Quantity        uint16
		Payload         []uint16
	}
)

func (hE *HistoryEvent) LogPrint() {
	log.Printf("\n\nSlave ID: %d\n", hE.Header.SlaveID)
	log.Printf("\n Transaction № %s", hE.Header.TransactionID)
	log.Println("\n Request:")
	hE.Handshake.Request.LogPrint()
	log.Println("\n Response:")
	hE.Handshake.Response.LogPrint()
	log.Printf("\n Transaction time: %v", hE.TransactionTime)
}

func (hdhk *Handshake) RequestUnmarshal(payload []byte) {
	if conf.WorkMode == "rtu_over_tcp" {
		functionID := payload[1]
		if slices.Contains([]byte{
			byte(conf.Functions.CoilsRead),
			byte(conf.Functions.DIRead),
			byte(conf.Functions.HRRead),
			byte(conf.Functions.IRRead),
			byte(conf.Functions.CoilsSimpleWrite),
			byte(conf.Functions.HRSimpleWrite)}, functionID) {
			hdhk.Request = new(RTUOverTCPRequest123456Response56)
		} else if slices.Contains([]byte{byte(conf.Functions.CoilsMultipleWrite), byte(conf.Functions.HRMultipleWrite)}, functionID) {
			hdhk.Request = new(RTUOverTCPMultipleWriteRequest)
		}
	} else {
		hdhk.Request = new(TCPRequest)
	}
	hdhk.Request.Unmarshal(payload)
}

func (hdhk *Handshake) ResponseUnmarshal(payload []byte) {
	if conf.WorkMode == "rtu_over_tcp" {
		functionID := payload[1]
		if slices.Contains([]byte{
			byte(conf.Functions.CoilsRead),
			byte(conf.Functions.DIRead),
			byte(conf.Functions.HRRead),
			byte(conf.Functions.IRRead)}, functionID) {
			hdhk.Response = new(RTUOverTCPReadResponse)
		} else if slices.Contains([]byte{byte(conf.Functions.CoilsSimpleWrite), byte(conf.Functions.HRSimpleWrite)}, functionID) {
			hdhk.Response = new(RTUOverTCPRequest123456Response56)
		} else if slices.Contains([]byte{byte(conf.Functions.CoilsMultipleWrite), byte(conf.Functions.HRMultipleWrite)}, functionID) {
			hdhk.Response = new(RTUOverTCPMultipleWriteResponse)
		} else {
			hdhk.Response = new(RTUOverTCPErrorResponse)
		}
	} else {
		hdhk.Response = new(TCPResponse)
	}
	hdhk.Response.Unmarshal(payload)
}

func (hdhk *Handshake) Marshal() (data EmulationData, err error) {
	data.FunctionID = hdhk.Response.GetFunctionID()
	data.IsReadOperation = !slices.Contains([]uint16{
		conf.Functions.CoilsSimpleWrite,
		conf.Functions.HRSimpleWrite,
		conf.Functions.CoilsMultipleWrite,
		conf.Functions.HRMultipleWrite}, data.FunctionID)
	address := hdhk.Request.MarshalAddress()
	data.Address = address[0] + address[1]
	quantity := hdhk.Request.MarshalQuantity()
	data.Quantity = quantity[0] + quantity[1]
	if data.IsReadOperation {
		if data.Payload, err = hdhk.Response.MarshalPayload(); err != nil {
			err = fmt.Errorf("error marshaling current handshake: %s", err)
			return
		}
		if len(data.Payload) != int(data.Quantity) {
			for {
				if len(data.Payload) == int(data.Quantity) {
					break
				}
				data.Payload = append(data.Payload, 0)
			}
		}
	} else {
		if data.Payload, err = hdhk.Request.MarshalPayload(); err != nil {
			err = fmt.Errorf("error marshaling current handshake: %s", err)
			return
		}
	}
	return
}

func (hdhk *Handshake) TransactionErrorCheck() bool {
	return hdhk.Response.GetFunctionID()>>7 == 0b1
}

func InputsPayloadPreprocessing[T uint16 | byte](data []T) (payload []uint16, err error) {
	for _, currentByte := range data {
		currentBinaryByte := strings.Split(strconv.FormatUint(uint64(currentByte), 2), "")
		for currentIndex := len(currentBinaryByte) - 1; currentIndex > -1; currentIndex-- {
			var currentIntBuffer int
			if currentIntBuffer, err = strconv.Atoi(currentBinaryByte[currentIndex]); err != nil {
				err = fmt.Errorf("error on marshaling binary read data: %s", err)
				return
			}
			payload = append(payload, uint16(currentIntBuffer))
		}
	}
	return
}

func RegistersPayloadPreprocessing[T uint16 | byte](data []T) (payload []uint16, err error) {
	for currentIndex := 0; currentIndex < len(data); currentIndex += 2 {
		var currentByte uint64
		if currentByte, err = strconv.ParseUint(fmt.Sprintf("%s%s",
			strconv.FormatUint(uint64(data[currentIndex]), 2), strconv.FormatUint(uint64(data[currentIndex+1]), 2)), 2, 64); err != nil {
			err = fmt.Errorf("error on marshaling registers data: %s", err)
			return
		}
		payload = append(payload, uint16(currentByte))
	}
	return
}
