package structs

import (
	"log"
	"slices"
)

type (
	DataPayload interface {
		MarshalPayload() []byte
		MarshalCheck() []byte
		Unmarshal([]byte)
		LogPrint()
	}
	TCPPacket interface {
		UnmarshalHeader([]byte)
		MarshalData() MarshaledData
		UnmarshalData([]byte)
		GetHeader() MBAPHeader
		LogPrint()
	}
	MBAPHeader struct {
		TransactionID []byte // [hight leve, low level]
		Protocol      string
		BodyLength    byte
		UnitID        byte
		FunctionType  byte
	}
	TCPPacketRequest struct {
		Header       MBAPHeader
		AddressStart []byte
		Data         DataPayload
	}
	TCPPacketResponse struct {
		Header MBAPHeader
		Data   DataPayload
	}
	ReadRequest struct {
		NumberReadingBits []byte // number of reading bits
	}
	ReadBitResponse struct { // for coils and DI
		NumberBits byte
		Bits       byte // like: [0, 1]
	}
	ReadByteResponse struct { // for HR and IR
		NumberBits byte
		Data       []byte // like: [0, 26, 0, 130]; len = numberBytes
	}
	WriteSimpleRequest struct {
		Payload []byte // like: [0, 6]
	}
	WriteMultipleRequest struct {
		NumberRegisters []byte // 2 bits, like: [0, 3] or [0, 2]
		NumberBits      byte
		Data            []byte // like: [0, 45, 0, 35]; len = numberRegisters[1]
	}
	WriteSimpleResponse struct {
		AddressStart []byte
		WrittenBits  []byte
	}
	WriteMultipleResponse struct {
		AddressStart           []byte
		NumberWrittenRegisters []byte
	}
	MarshaledData struct {
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

func (pReq *TCPPacketRequest) UnmarshalHeader(payload []byte) {
	if len(payload) < 10 {
		log.Println("Error: insufficient payload length")
		return
	}
	pReq.Header.Unmarshal(payload)
	pReq.AddressStart = payload[8:10]
}

func (pReq *TCPPacketRequest) MarshalData() (data MarshaledData) {
	data.AddressStart = pReq.AddressStart
	data.CheckField = pReq.Data.MarshalCheck()
	data.Payload = pReq.Data.MarshalPayload()
	return
}

func (pReq *TCPPacketRequest) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{1, 2, 3, 4}, pReq.Header.FunctionType) {
		pReq.Data = new(ReadRequest)
	} else if slices.Contains([]byte{5, 6}, pReq.Header.FunctionType) {
		pReq.Data = new(WriteSimpleRequest)
	} else if slices.Contains([]byte{15, 16}, pReq.Header.FunctionType) {
		pReq.Data = new(WriteMultipleRequest)
	}
	pReq.Data.Unmarshal(payload)
}

func (pReq *TCPPacketRequest) GetHeader() MBAPHeader {
	return pReq.Header
}

func (pReq *TCPPacketRequest) LogPrint() {
	log.Println("  Header:")
	pReq.Header.LogPrint()
	log.Printf("  Address Start: %v", pReq.AddressStart)
	log.Println("  Data:")
	pReq.Data.LogPrint()
}

func (pRes *TCPPacketResponse) UnmarshalHeader(payload []byte) {
	pRes.Header.Unmarshal(payload)
}

func (pRes *TCPPacketResponse) MarshalData() (data MarshaledData) {
	data.AddressStart = nil
	data.CheckField = pRes.Data.MarshalCheck()
	data.Payload = pRes.Data.MarshalPayload()
	return
}

func (pRes *TCPPacketResponse) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{1, 2}, pRes.Header.FunctionType) {
		pRes.Data = new(ReadBitResponse)
	} else if slices.Contains([]byte{3, 4}, pRes.Header.FunctionType) {
		pRes.Data = new(ReadByteResponse)
	} else if slices.Contains([]byte{5, 6}, pRes.Header.FunctionType) {
		pRes.Data = new(WriteSimpleResponse)
	} else if slices.Contains([]byte{15, 16}, pRes.Header.FunctionType) {
		pRes.Data = new(WriteMultipleResponse)
	}
	pRes.Data.Unmarshal(payload)
}

func (pRes *TCPPacketResponse) GetHeader() MBAPHeader {
	return pRes.Header
}

func (pRes *TCPPacketResponse) LogPrint() {
	log.Println("  Header:")
	pRes.Header.LogPrint()
	log.Println("  Data:")
	pRes.Data.LogPrint()
}

func (rReq *ReadRequest) MarshalPayload() []byte {
	return rReq.NumberReadingBits
}

func (rReq *ReadRequest) MarshalCheck() []byte {
	return rReq.MarshalPayload()
}

func (rReq *ReadRequest) Unmarshal(payload []byte) {
	if len(payload) < 11 {
		log.Println("Error: insufficient payload length")
		return
	}
	rReq.NumberReadingBits = payload[10:]
}

func (rReq *ReadRequest) LogPrint() {
	log.Printf("   Number reading bits: %v\n", rReq.NumberReadingBits)
}

func (rBiRes *ReadBitResponse) MarshalPayload() (payload []byte) {
	return []byte{rBiRes.Bits}
}

func (rBiRes *ReadBitResponse) MarshalCheck() []byte {
	return []byte{rBiRes.NumberBits / 2}
}

func (rBiRes *ReadBitResponse) Unmarshal(payload []byte) {
	rBiRes.NumberBits = payload[8]
	rBiRes.Bits = payload[9]
}

func (rBiRes *ReadBitResponse) LogPrint() {
	log.Printf("   Count response bit: %v\n", rBiRes.NumberBits)
	log.Printf("   Response bit: %v\n", rBiRes.Bits)
}

func (rByRes *ReadByteResponse) MarshalPayload() []byte {
	return rByRes.Data
}

func (rByRes *ReadByteResponse) MarshalCheck() []byte {
	return []byte{rByRes.NumberBits / 2}
}

func (rByRes *ReadByteResponse) Unmarshal(payload []byte) {
	rByRes.NumberBits = payload[8]
	rByRes.Data = payload[9:]
}

func (rByRes *ReadByteResponse) LogPrint() {
	log.Printf("   Number of respoce bits: %v\n", rByRes.NumberBits)
	log.Printf("   Response bits: %v\n", rByRes.Data)
}

func (wSReq *WriteSimpleRequest) MarshalPayload() (payload []byte) {
	return wSReq.Payload
}

func (wReq *WriteSimpleRequest) MarshalCheck() []byte {
	return wReq.MarshalPayload()
}

func (wSReq *WriteSimpleRequest) Unmarshal(payload []byte) {
	if len(payload) < 11 {
		log.Println("Error: insufficient payload length")
		return
	}
	wSReq.Payload = payload[10:]
}

func (wSReq *WriteSimpleRequest) LogPrint() {
	log.Printf("   Writing bits: %v\n", wSReq.Payload)
}

func (wMReq *WriteMultipleRequest) MarshalPayload() []byte {
	return wMReq.Data
}

func (wMReq *WriteMultipleRequest) MarshalCheck() []byte {
	return wMReq.NumberRegisters
}

func (wMReq *WriteMultipleRequest) Unmarshal(payload []byte) {
	if len(payload) < 14 {
		log.Println("Error: insufficient payload length")
		return
	}
	wMReq.NumberRegisters = payload[10:12]
	wMReq.NumberBits = payload[12]
	wMReq.Data = payload[13:]
}

func (wMReq *WriteMultipleRequest) LogPrint() {
	log.Printf("   Number of writting registers: %v\n", wMReq.NumberRegisters)
	log.Printf("   Number of writting bits: %v\n", wMReq.NumberBits)
	log.Printf("   Writting bits: %v\n", wMReq.Data)
}

func (wSRes *WriteSimpleResponse) MarshalPayload() (payload []byte) {
	return wSRes.WrittenBits
}
func (wSRes *WriteSimpleResponse) MarshalCheck() []byte {
	return wSRes.MarshalPayload()
}

func (wSRes *WriteSimpleResponse) Unmarshal(payload []byte) {
	wSRes.AddressStart = payload[8:10]
	wSRes.WrittenBits = payload[10:]
}

func (wSRes *WriteSimpleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wSRes.AddressStart)
	log.Printf("   Written bits: %v\n", wSRes.WrittenBits)
}

func (wMRes *WriteMultipleResponse) MarshalPayload() (payload []byte) {
	return wMRes.NumberWrittenRegisters
}

func (wMRes *WriteMultipleResponse) MarshalCheck() []byte {
	return wMRes.MarshalPayload()
}

func (wMRes *WriteMultipleResponse) Unmarshal(payload []byte) {
	wMRes.AddressStart = payload[8:10]
	wMRes.NumberWrittenRegisters = payload[10:]
}

func (wMRes *WriteMultipleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wMRes.AddressStart)
	log.Printf("   Number written registers: %v\n", wMRes.NumberWrittenRegisters)
}
