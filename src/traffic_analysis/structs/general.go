package structs

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"slices"
	"sort"
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

func (hdhk *Handshake) RequestUnmarshal(workMode string, payload []byte) {
	switch workMode {
	case conf.Protocols.RTUOverTCP:
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
	case conf.Protocols.TCP:
		hdhk.Request = new(TCPRequest)
	default:
		log.Fatal("Error on unmarshaling request: invalid protocol")
	}
	hdhk.Request.Unmarshal(payload)
}

func (hdhk *Handshake) ResponseUnmarshal(workMode string, payload []byte) {
	switch workMode {
	case conf.Protocols.RTUOverTCP:
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
	case conf.Protocols.TCP:
		hdhk.Response = new(TCPResponse)
	default:
		log.Fatal("Error on unmarshaling request: invalid protocol")
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
	if data.Address, err = BytesToDecimal(hdhk.Request.MarshalAddress()); err != nil {
		err = fmt.Errorf("error on marshaliing emulation data address: %s", err)
		return
	}
	if data.Quantity, err = BytesToDecimal(hdhk.Request.MarshalQuantity()); err != nil {
		err = fmt.Errorf("error on marshaliing emulation data quantity: %s", err)
		return
	}
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

func (sH *ServerHistory) SelfClean() {
	deleteIndices := []int{}
	for currentIndex, currentHistoryEvent := range sH.Transactions {
		if currentHistoryEvent.Handshake.Request == nil || currentHistoryEvent.Handshake.Response == nil {
			deleteIndices = append(deleteIndices, currentIndex)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(deleteIndices)))
	for _, currentDeleteIndex := range deleteIndices {
		sH.Transactions = append(sH.Transactions[:currentDeleteIndex], sH.Transactions[currentDeleteIndex+1:]...)
	}
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
		var currentByte uint16
		if currentByte, err = BytesToDecimal(data[currentIndex : currentIndex+2]); err != nil {
			err = fmt.Errorf("error on marshaling registers data: %s", err)
			return
		}
		payload = append(payload, uint16(currentByte))
	}
	return
}

func BytesToDecimal[T uint16 | byte](bytes []T) (result uint16, err error) {
	var hexBuffer string
	for _, curretnByte := range bytes {
		hexBuffer = fmt.Sprintf("%s%s", hexBuffer, strconv.FormatUint(uint64(curretnByte), 16))
	}
	var resultBuffer uint64
	if resultBuffer, err = strconv.ParseUint(hexBuffer, 16, 64); err != nil {
		err = fmt.Errorf("error on parsing bytes: %s", err)
		return
	}
	result = uint16(resultBuffer)
	return
}
