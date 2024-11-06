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
	// WriteSimpleRequest struct {
	// 	payload []byte // like: [0, 6]
	// }
	// WriteMultipleRequest struct {
	// 	payload struct {
	// 		numberRegisters []byte // 2 bits, like: [0, 3] or [0, 2]
	// 		numberBits      byte
	// 		data            [][]byte // like: [[0, 45], [0, 35]]; len = numberRegisters[1]
	// 	}
	// }
	// WriteSingleResponce struct {
	// 	payload []byte // written bits
	// }
	// WriteMultipleResponce struct {
	// 	payload byte // number of written bits
	// }
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
	var payloadData [][]byte
	var workData []byte
	for currentIndex, currentBit := range payload[9:] {
		if currentIndex%2 == 0 {
			workData = []byte{currentBit}
		} else {
			workData = append(workData, currentBit)
			payloadData = append(payloadData, workData)
		}
	}
	rByRes.numberBits = payload[8]
	rByRes.data = payloadData
}

func (rByRes *ReadByteResponce) LogPrint() {
	log.Printf("   Number of respoce bits: %v\n", rByRes.numberBits)
	log.Printf("   Responce bits: %v\n", rByRes.data)
}
