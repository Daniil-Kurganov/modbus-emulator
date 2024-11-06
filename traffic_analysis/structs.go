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
		Responce TCPPacket
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
	TCPPacketResponce struct {
		Header MBAPHeader
		Data   DataPayload
	}
	ReadRequest struct {
		numberReadingBits []byte // number of reading bits
	}
	ReadBitResponce struct { // for coils and DI
		bits []byte // like: [0, 1]
	}
	ReadByteResponce struct { // for HR and IR
		numberBits byte
		data       [][]byte // like: [[0, 26], [0, 130]]; len = numberBytes
	}
	WriteSimpleRequest struct {
		payload []byte // like: [0, 6]
	}
	WriteMultipleRequest struct {
		numberRegisters []byte // 2 bits, like: [0, 3] or [0, 2]
		numberBits      byte
		data            []byte // like: [[0, 45], [0, 35]]; len = numberRegisters[1]
	}
	WriteSimpleResponce struct {
		addressStart []byte
		writtenBits  []byte
	}
	WriteMultipleResponce struct {
		addressStart           []byte
		numberWrittenRegisters []byte
	}
)

func (h *MBAPHeader) Unmarshal(payload []byte) {
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

func (pRes *TCPPacketResponce) UnmarshalHeader(payload []byte) {
	pRes.Header.Unmarshal(payload)
}

func (pRes *TCPPacketResponce) UnmarshalData(payload []byte) {
	if slices.Contains([]byte{1, 2}, pRes.Header.FunctionType) {
		pRes.Data = new(ReadBitResponce)
	} else if slices.Contains([]byte{3, 4}, pRes.Header.FunctionType) {
		pRes.Data = new(ReadByteResponce)
	} else if slices.Contains([]byte{5, 6}, pRes.Header.FunctionType) {
		pRes.Data = new(WriteSimpleResponce)
	} else if slices.Contains([]byte{15, 16}, pRes.Header.FunctionType) {
		pRes.Data = new(WriteMultipleResponce)
	}
	pRes.Data.Unmarshal(payload)
}

func (pRes *TCPPacketResponce) LogPrint() {
	log.Println("  Header:")
	pRes.Header.LogPrint()
	log.Println("  Data:")
	pRes.Data.LogPrint()
}

func (rReq *ReadRequest) Marshal() (payload []byte) {
	return
}

func (rReq *ReadRequest) Unmarshal(payload []byte) {
	rReq.numberReadingBits = payload[10:]
}

func (rReq *ReadRequest) LogPrint() {
	log.Printf("   Number reading bits: %v\n", rReq.numberReadingBits)
}

func (rBiRes *ReadBitResponce) Marshal() (payload []byte) {
	return
}

func (rBiRes *ReadBitResponce) Unmarshal(payload []byte) {
	rBiRes.bits = payload[8:]
}

func (rBiRes *ReadBitResponce) LogPrint() {
	log.Printf("   Responce bit: %v\n", rBiRes.bits)
}

func (rByRes *ReadByteResponce) Marshal() (payload []byte) {
	return
}

func (rByRes *ReadByteResponce) Unmarshal(payload []byte) {
	rByRes.numberBits = payload[8]
	rByRes.data = convertPayloadData(payload[9:])
}

func (rByRes *ReadByteResponce) LogPrint() {
	log.Printf("   Number of respoce bits: %v\n", rByRes.numberBits)
	log.Printf("   Responce bits: %v\n", rByRes.data)
}

func (wSReq *WriteSimpleRequest) Marshal() (payload []byte) {
	return
}

func (wSReq *WriteSimpleRequest) Unmarshal(payload []byte) {
	wSReq.payload = payload[10:]
}

func (wSReq *WriteSimpleRequest) LogPrint() {
	log.Printf("   Writing bits: %v\n", wSReq.payload)
}

func (wMReq *WriteMultipleRequest) Marshal() (payload []byte) {
	return
}

func (wMReq *WriteMultipleRequest) Unmarshal(payload []byte) {
	wMReq.numberRegisters = payload[10:12]
	wMReq.numberBits = payload[12]
	wMReq.data = payload[13:]
}

func (wMReq *WriteMultipleRequest) LogPrint() {
	log.Printf("   Number of writting registers: %v\n", wMReq.numberRegisters)
	log.Printf("   Number of writting bits: %v\n", wMReq.numberBits)
	log.Printf("   Writting bits: %v\n", wMReq.data)
}

func (wSRes *WriteSimpleResponce) Marshal() (payload []byte) {
	return
}

func (wSRes *WriteSimpleResponce) Unmarshal(payload []byte) {
	wSRes.addressStart = payload[8:10]
	wSRes.writtenBits = payload[10:]
}

func (wSRes *WriteSimpleResponce) LogPrint() {
	log.Printf("   Address start: %v\n", wSRes.addressStart)
	log.Printf("   Written bits: %v\n", wSRes.writtenBits)
}

func (wMRes *WriteMultipleResponce) Marshal() (payload []byte) {
	return
}

func (wMRes *WriteMultipleResponce) Unmarshal(payload []byte) {
	wMRes.addressStart = payload[8:10]
	wMRes.numberWrittenRegisters = payload[10:]
}

func (wMRes *WriteMultipleResponce) LogPrint() {
	log.Printf("   Address start: %v\n", wMRes.addressStart)
	log.Printf("   Number written registers: %v\n", wMRes.numberWrittenRegisters)
}

func convertPayloadData(data []byte) (payloadData [][]byte) {
	var workData []byte
	for currentIndex, currentBit := range data {
		if currentIndex%2 == 0 {
			workData = []byte{currentBit}
		} else {
			workData = append(workData, currentBit)
			payloadData = append(payloadData, workData)
		}
	}
	return
}
