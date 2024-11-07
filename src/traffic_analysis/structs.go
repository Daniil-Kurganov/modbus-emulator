package trafficanalysis

import (
	"log"
	"slices"
)

type (
	DataPayload interface {
		Marshal() []byte
		Unmarshal([]byte)
		LogPrint()
	}
	TCPPacket interface {
		UnmarshalHeader([]byte)
		UnmarshalData([]byte)
		LogPrint()
	}
	Handshake struct {
		Request  TCPPacket
		Response TCPPacket
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
		numberBits byte
		data       [][]byte // like: [[0, 26], [0, 130]]; len = numberBytes
	}
	WriteSimpleRequest struct {
		Payload []byte // like: [0, 6]
	}
	WriteMultipleRequest struct {
		NumberRegisters []byte // 2 bits, like: [0, 3] or [0, 2]
		NumberBits      byte
		Data            []byte // like: [[0, 45], [0, 35]]; len = numberRegisters[1]
	}
	WriteSimpleResponse struct {
		AddressStart []byte
		WrittenBits  []byte
	}
	WriteMultipleResponse struct {
		AddressStart           []byte
		NumberWrittenRegisters []byte
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

func (pRes *TCPPacketResponse) LogPrint() {
	log.Println("  Header:")
	pRes.Header.LogPrint()
	log.Println("  Data:")
	pRes.Data.LogPrint()
}

func (rReq *ReadRequest) Marshal() (payload []byte) {
	return
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

func (rBiRes *ReadBitResponse) Marshal() (payload []byte) {
	return
}

func (rBiRes *ReadBitResponse) Unmarshal(payload []byte) {
	rBiRes.NumberBits = payload[8]
	rBiRes.Bits = payload[9]
}

func (rBiRes *ReadBitResponse) LogPrint() {
	log.Printf("   Count response bit: %v\n", rBiRes.NumberBits)
	log.Printf("   Response bit: %v\n", rBiRes.Bits)
}

func (rByRes *ReadByteResponse) Marshal() (payload []byte) {
	return
}

func (rByRes *ReadByteResponse) Unmarshal(payload []byte) {
	rByRes.numberBits = payload[8]
	var workData []byte
	for currentIndex, currentBit := range payload[9:] {
		if currentIndex%2 == 0 {
			workData = []byte{currentBit}
		} else {
			workData = append(workData, currentBit)
			rByRes.data = append(rByRes.data, workData)
		}
	}
}

func (rByRes *ReadByteResponse) LogPrint() {
	log.Printf("   Number of respoce bits: %v\n", rByRes.numberBits)
	log.Printf("   Response bits: %v\n", rByRes.data)
}

func (wSReq *WriteSimpleRequest) Marshal() (payload []byte) {
	return
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

func (wMReq *WriteMultipleRequest) Marshal() (payload []byte) {
	return
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

func (wSRes *WriteSimpleResponse) Marshal() (payload []byte) {
	return
}

func (wSRes *WriteSimpleResponse) Unmarshal(payload []byte) {
	wSRes.AddressStart = payload[8:10]
	wSRes.WrittenBits = payload[10:]
}

func (wSRes *WriteSimpleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wSRes.AddressStart)
	log.Printf("   Written bits: %v\n", wSRes.WrittenBits)
}

func (wMRes *WriteMultipleResponse) Marshal() (payload []byte) {
	return
}

func (wMRes *WriteMultipleResponse) Unmarshal(payload []byte) {
	wMRes.AddressStart = payload[8:10]
	wMRes.NumberWrittenRegisters = payload[10:]
}

func (wMRes *WriteMultipleResponse) LogPrint() {
	log.Printf("   Address start: %v\n", wMRes.AddressStart)
	log.Printf("   Number written registers: %v\n", wMRes.NumberWrittenRegisters)
}
