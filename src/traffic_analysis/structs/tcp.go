package structs

import (
	"log"
	"slices"
)

type (
	DataPayload interface {
		GetQuantityRegisters() []uint16
		MarshalPayload() []byte
		MarshalCheck() []byte
		Unmarshal([]byte)
		LogPrint()
	}
	// TCPPacket interface {
	// 	UnmarshalHeader([]byte)
	// 	MarshalData() MarshaledData
	// 	UnmarshalData([]byte)
	// 	GetHeader() MBAPHeader
	// 	LogPrint()
	// }
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
	for _, currentByte := range pReq.Data.MarshalPayload() {
		payload = append(payload, uint16(currentByte))
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

// func (pReq *TCPRequest) MarshalData() (data TCPMarshaledData) {
// 	data.AddressStart = pReq.AddressStart
// 	data.CheckField = pReq.Data.MarshalCheck()
// 	data.Payload = pReq.Data.MarshalPayload()
// 	return
// }

func (pReq *TCPRequest) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{1, 2, 3, 4}, pReq.Header.FunctionType) {
		pReq.Data = new(TCPReadRequest)
	} else if slices.Contains([]byte{5, 6}, pReq.Header.FunctionType) {
		pReq.Data = new(TCPWriteSimpleRequest)
	} else if slices.Contains([]byte{15, 16}, pReq.Header.FunctionType) {
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
	for _, currentByte := range pRes.Data.MarshalPayload() {
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

func (pRes *TCPResponse) MarshalData() (data TCPMarshaledData) {
	data.AddressStart = nil
	data.CheckField = pRes.Data.MarshalCheck()
	data.Payload = pRes.Data.MarshalPayload()
	return
}

func (pRes *TCPResponse) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{1, 2}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPReadBitResponse)
	} else if slices.Contains([]byte{3, 4}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPReadByteResponse)
	} else if slices.Contains([]byte{5, 6}, pRes.Header.FunctionType) {
		pRes.Data = new(TCPWriteSimpleResponse)
	} else if slices.Contains([]byte{15, 16}, pRes.Header.FunctionType) {
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

func (rReq *TCPReadRequest) MarshalPayload() []byte {
	return rReq.NumberReadingBits
}

func (rReq *TCPReadRequest) MarshalCheck() []byte {
	return rReq.MarshalPayload()
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

func (rBiRes *TCPReadBitResponse) MarshalPayload() (payload []byte) {
	return []byte{rBiRes.Bits}
}

func (rBiRes *TCPReadBitResponse) MarshalCheck() []byte {
	return []byte{rBiRes.NumberBits / 2}
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

func (rByRes *TCPReadByteResponse) MarshalPayload() []byte {
	return rByRes.Data
}

func (rByRes *TCPReadByteResponse) MarshalCheck() []byte {
	return []byte{rByRes.NumberBits / 2}
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

func (wSReq *TCPWriteSimpleRequest) MarshalPayload() (payload []byte) {
	return wSReq.Payload
}

func (wReq *TCPWriteSimpleRequest) MarshalCheck() []byte {
	return wReq.MarshalPayload()
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

func (wMReq *TCPWriteMultipleRequest) MarshalPayload() []byte {
	return wMReq.Data
}

func (wMReq *TCPWriteMultipleRequest) MarshalCheck() []byte {
	return wMReq.NumberRegisters
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

func (wSRes *TCPWriteSimpleResponse) MarshalPayload() (payload []byte) {
	return wSRes.WrittenBits
}
func (wSRes *TCPWriteSimpleResponse) MarshalCheck() []byte {
	return wSRes.MarshalPayload()
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

func (wMRes *TCPWriteMultipleResponse) MarshalPayload() (payload []byte) {
	return wMRes.NumberWrittenRegisters
}

func (wMRes *TCPWriteMultipleResponse) MarshalCheck() []byte {
	return wMRes.MarshalPayload()
}

func (wMRes *TCPWriteMultipleResponse) Unmarshal(payload []byte) {
	wMRes.AddressStart = payload[8:10]
	wMRes.NumberWrittenRegisters = payload[10:]
}

func (wMRes *TCPWriteMultipleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wMRes.AddressStart)
	log.Printf("   Number written registers: %v\n", wMRes.NumberWrittenRegisters)
}
