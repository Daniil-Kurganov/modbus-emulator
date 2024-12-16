package structs

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
	"reflect"
	"slices"
)

type (
	DataPayload interface {
		GetQuantityRegisters() []uint16
		MarshalPayload() ([]uint16, error)
		Unmarshal([]byte)
		LogPrint()
	}

	MBAPHeader struct {
		TransactionID []byte // [hight leve, low level]
		Protocol      string
		BodyLength    byte
		UnitID        byte
		FunctionType  byte
	}
	TCPRequest struct {
		Header       MBAPHeader
		AddressStart []byte
		Data         DataPayload
	}
	TCPResponse struct {
		Header MBAPHeader
		Data   DataPayload
	}
	TCPReadRequest struct {
		NumberReadingBits []byte // number of reading bits
	}
	TCPReadBitResponse struct { // for coils and DI
		NumberBits byte
		Bits       byte // like: [0, 1]
	}
	TCPReadByteResponse struct { // for HR and IR
		NumberBits byte
		Data       []byte // like: [0, 26, 0, 130]; len = numberBytes
	}
	TCPWriteSimpleRequest struct {
		Payload []byte // like: [0, 6]
	}
	TCPWriteMultipleRequest struct {
		NumberRegisters []byte // 2 bits, like: [0, 3] or [0, 2]
		NumberBits      byte
		Data            []byte // like: [0, 45, 0, 35]; len = numberRegisters[1]
	}
	TCPWriteSimpleResponse struct {
		AddressStart []byte
		WrittenBits  []byte
	}
	TCPWriteMultipleResponse struct {
		AddressStart           []byte
		NumberWrittenRegisters []byte
	}
	TCPMarshaledData struct {
		AddressStart []byte
		CheckField   []byte
		Payload      []byte
	}
)

func (h *MBAPHeader) Unmarshal(payload []byte) {
	if len(payload) < 8 {
		log.Println("Error: insufficient payload length")
		return
	}
	h.TransactionID = payload[:2]
	if payload[2] == 0 && payload[3] == 0 {
		h.Protocol = "modbus"
	} else {
		h.Protocol = "unknown"
		return
	}
	h.BodyLength = payload[4] + payload[5]
	h.UnitID = payload[6]
	h.FunctionType = payload[7]
}

func (h *MBAPHeader) LogPrint() {
	log.Printf("   Transaction ID: %v\n", h.TransactionID)
	log.Printf("   Protocol: %s\n", h.Protocol)
	log.Printf("   Body length: %v\n", h.BodyLength)
	log.Printf("   Unit ID: %v\n", h.UnitID)
	log.Printf("   Function type: %v\n", h.FunctionType)
}

func (pReq *TCPRequest) Unmarshal(payload []byte) {
	pReq.UnmarshalHeader(payload)
	pReq.UnmarshalData(payload)
}

func (pReq *TCPRequest) MarshalPayload() (payload []uint16, err error) {
	if payload, err = pReq.Data.MarshalPayload(); err != nil {
		err = fmt.Errorf("error on marshaling request: %s", err)
		return
	}
	if pReq.Header.FunctionType == byte(conf.Functions.HRMultipleWrite) {
		if payload, err = RegistersPayloadPreprocessing(payload); err != nil {
			err = fmt.Errorf("error on marshaling request: %s", err)
			return
		}
	}
	return
}

func (pReq *TCPRequest) LogPrint() {
	log.Println("  Header:")
	pReq.Header.LogPrint()
	log.Printf("  Address Start: %v", pReq.AddressStart)
	log.Println("  Data:")
	pReq.Data.LogPrint()
}

func (pReq *TCPRequest) MarshalAddress() (address []uint16) {
	for _, currentAddressIndex := range pReq.AddressStart {
		address = append(address, uint16(currentAddressIndex))
	}
	return
}

func (pReq *TCPRequest) MarshalQuantity() []uint16 {
	return pReq.Data.GetQuantityRegisters()
}

func (pReq *TCPRequest) UnmarshalHeader(payload []byte) {
	if len(payload) < 10 {
		log.Println("Error: insufficient payload length")
		return
	}
	pReq.Header.Unmarshal(payload)
	pReq.AddressStart = payload[8:10]
}

func (pReq *TCPRequest) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{
		byte(conf.Functions.CoilsRead),
		byte(conf.Functions.DIRead),
		byte(conf.Functions.HRRead),
		byte(conf.Functions.IRRead)}, pReq.Header.FunctionType) {
		pReq.Data = new(TCPReadRequest)
	} else if slices.Contains([]byte{byte(conf.Functions.CoilsSimpleWrite), byte(conf.Functions.HRSimpleWrite)}, pReq.Header.FunctionType) {
		pReq.Data = new(TCPWriteSimpleRequest)
	} else if slices.Contains([]byte{byte(conf.Functions.CoilsMultipleWrite), byte(conf.Functions.HRMultipleWrite)}, pReq.Header.FunctionType) {
		pReq.Data = new(TCPWriteMultipleRequest)
	}
	pReq.Data.Unmarshal(payload)
}

func (pReq *TCPRequest) GetHeader() MBAPHeader {
	return pReq.Header
}

func (pRes *TCPResponse) Unmarshal(payload []byte) {
	pRes.UnmarshalHeader(payload)
	pRes.UnmarshalData(payload)
}

func (pRes *TCPResponse) MarshalPayload() (payload []uint16, err error) {
	var workPayload []uint16
	if workPayload, err = pRes.Data.MarshalPayload(); err != nil {
		err = fmt.Errorf("error on marshaling response: %s", err)
		return
	}
	for _, currentByte := range workPayload {
		payload = append(payload, uint16(currentByte))
	}
	return
}

func (pRes *TCPResponse) LogPrint() {
	log.Println("  Header:")
	pRes.Header.LogPrint()
	log.Println("  Data:")
	pRes.Data.LogPrint()
}

func (pRes *TCPResponse) GetFunctionID() uint16 {
	return uint16(pRes.Header.FunctionType)
}

func (pRes *TCPResponse) UnmarshalHeader(payload []byte) {
	pRes.Header.Unmarshal(payload)
}

func (pRes *TCPResponse) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{byte(conf.Functions.CoilsRead), byte(conf.Functions.DIRead)}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPReadBitResponse)
	} else if slices.Contains([]byte{byte(conf.Functions.HRRead), byte(conf.Functions.IRRead)}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPReadByteResponse)
	} else if slices.Contains([]byte{byte(conf.Functions.CoilsSimpleWrite), byte(conf.Functions.HRSimpleWrite)}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPWriteSimpleResponse)
	} else if slices.Contains([]byte{byte(conf.Functions.CoilsMultipleWrite), byte(conf.Functions.HRMultipleWrite)}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPWriteMultipleResponse)
	}
	pRes.Data.Unmarshal(payload)
}

func (pRes *TCPResponse) GetHeader() MBAPHeader {
	return pRes.Header
}

func (rReq *TCPReadRequest) GetQuantityRegisters() (quantity []uint16) {
	for _, currentQuantityLevel := range rReq.NumberReadingBits {
		quantity = append(quantity, uint16(currentQuantityLevel))
	}
	return
}

func (rReq *TCPReadRequest) MarshalPayload() (payload []uint16, err error) {
	for _, currentByte := range rReq.NumberReadingBits {
		payload = append(payload, uint16(currentByte))
	}
	return
}

func (rReq *TCPReadRequest) Unmarshal(payload []byte) {
	if len(payload) < 11 {
		log.Println("Error: insufficient payload length")
		return
	}
	rReq.NumberReadingBits = payload[10:]
}

func (rReq *TCPReadRequest) LogPrint() {
	log.Printf("   Number reading bits: %v\n", rReq.NumberReadingBits)
}

func (rBiRes *TCPReadBitResponse) GetQuantityRegisters() []uint16 {
	return []uint16{}
}

func (rBiRes *TCPReadBitResponse) MarshalPayload() (payload []uint16, err error) {
	if payload, err = InputsPayloadPreprocessing([]byte{rBiRes.Bits}); err != nil {
		err = fmt.Errorf("error on marshaling read data: %s", err)
	}
	return
}

func (rBiRes *TCPReadBitResponse) Unmarshal(payload []byte) {
	rBiRes.NumberBits = payload[8]
	rBiRes.Bits = payload[9]
}

func (rBiRes *TCPReadBitResponse) LogPrint() {
	log.Printf("   Count response bit: %v\n", rBiRes.NumberBits)
	log.Printf("   Response bit: %v\n", rBiRes.Bits)
}

func (rByRes *TCPReadByteResponse) GetQuantityRegisters() []uint16 {
	return []uint16{}
}

func (rByRes *TCPReadByteResponse) MarshalPayload() (payload []uint16, err error) {
	if payload, err = RegistersPayloadPreprocessing(rByRes.Data); err != nil {
		err = fmt.Errorf("error on marshaling read register's data: %s", err)
		return
	}
	return
}

func (rByRes *TCPReadByteResponse) Unmarshal(payload []byte) {
	rByRes.NumberBits = payload[8]
	rByRes.Data = payload[9:]
}

func (rByRes *TCPReadByteResponse) LogPrint() {
	log.Printf("   Number of respoce bits: %v\n", rByRes.NumberBits)
	log.Printf("   Response bits: %v\n", rByRes.Data)
}

func (wSReq *TCPWriteSimpleRequest) GetQuantityRegisters() (quantity []uint16) {
	return []uint16{0, 1}
}

func (wSReq *TCPWriteSimpleRequest) MarshalPayload() (payload []uint16, err error) {
	if reflect.DeepEqual(wSReq.Payload, []byte{255, 0}) {
		payload = []uint16{1}
	} else if reflect.DeepEqual(wSReq.Payload, []byte{0, 0}) {
		payload = []uint16{0}
	} else {
		if payload, err = RegistersPayloadPreprocessing(wSReq.Payload); err != nil {
			err = fmt.Errorf("error on marshaling read data: %s", err)
		}
	}
	return
}

func (wSReq *TCPWriteSimpleRequest) Unmarshal(payload []byte) {
	if len(payload) < 11 {
		log.Println("Error: insufficient payload length")
		return
	}
	wSReq.Payload = payload[10:]
}

func (wSReq *TCPWriteSimpleRequest) LogPrint() {
	log.Printf("   Writing bits: %v\n", wSReq.Payload)
}

func (wMreq *TCPWriteMultipleRequest) GetQuantityRegisters() (quantity []uint16) {
	for _, currentQuantityLevel := range wMreq.NumberRegisters {
		quantity = append(quantity, uint16(currentQuantityLevel))
	}
	return
}

func (wMReq *TCPWriteMultipleRequest) MarshalPayload() (payload []uint16, err error) {
	for _, currentByte := range wMReq.Data {
		payload = append(payload, uint16(currentByte))
	}
	return
}

func (wMReq *TCPWriteMultipleRequest) Unmarshal(payload []byte) {
	if len(payload) < 14 {
		log.Println("Error: insufficient payload length")
		return
	}
	wMReq.NumberRegisters = payload[10:12]
	wMReq.NumberBits = payload[12]
	wMReq.Data = payload[13:]
}

func (wMReq *TCPWriteMultipleRequest) LogPrint() {
	log.Printf("   Number of writting registers: %v\n", wMReq.NumberRegisters)
	log.Printf("   Number of writting bits: %v\n", wMReq.NumberBits)
	log.Printf("   Writting bits: %v\n", wMReq.Data)
}

func (wSRes *TCPWriteSimpleResponse) GetQuantityRegisters() []uint16 {
	return []uint16{}
}

func (wSRes *TCPWriteSimpleResponse) MarshalPayload() (payload []uint16, err error) {
	for _, currentByte := range wSRes.WrittenBits {
		payload = append(payload, uint16(currentByte))
	}
	return
}

func (wSRes *TCPWriteSimpleResponse) Unmarshal(payload []byte) {
	wSRes.AddressStart = payload[8:10]
	wSRes.WrittenBits = payload[10:]
}

func (wSRes *TCPWriteSimpleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wSRes.AddressStart)
	log.Printf("   Written bits: %v\n", wSRes.WrittenBits)
}

func (wMRes *TCPWriteMultipleResponse) GetQuantityRegisters() (quantity []uint16) {
	for _, currentQuantityLevel := range wMRes.NumberWrittenRegisters {
		quantity = append(quantity, uint16(currentQuantityLevel))
	}
	return
}

func (wMRes *TCPWriteMultipleResponse) MarshalPayload() (payload []uint16, err error) {
	for _, currentByte := range wMRes.NumberWrittenRegisters {
		payload = append(payload, uint16(currentByte))
	}
	return
}

func (wMRes *TCPWriteMultipleResponse) Unmarshal(payload []byte) {
	wMRes.AddressStart = payload[8:10]
	wMRes.NumberWrittenRegisters = payload[10:]
}

func (wMRes *TCPWriteMultipleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wMRes.AddressStart)
	log.Printf("   Number written registers: %v\n", wMRes.NumberWrittenRegisters)
}
