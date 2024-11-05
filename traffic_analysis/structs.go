package trafficanalysis

type (
	DataPayload interface {
		Marshal() []byte
		Unmarshal([]byte)
	}
	TCPPacket interface {
		UnmarshalHeader([]byte) MBAPHeader
	}
	Handshake struct {
		request  TCPPacketRequest
		responce TCPPacketResponce
	}
	MBAPHeader struct {
		TransactionID []byte // [hight leve, low level]
		Protocol      string
		BodyLength    byte
		UnitID        byte
		FunctionType  byte
		Data          DataPayload
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

func (p *TCPPacketRequest) UnmarshalHeader(payload []byte) MBAPHeader {
	p.Header.Unmarshal(payload)
	p.AddressStart = payload[8:10]
	return p.Header
}

func (p *TCPPacketResponce) UnmarshalHeader(payload []byte) MBAPHeader {
	p.Header.Unmarshal(payload)
	return p.Header
}

func (rReq *ReadRequest) Marshal() (payload []byte) {
	return
}

func (rReq *ReadRequest) Unmarshal(payload []byte) {
	rReq.numberReadingBits = payload[10:]
}

func (rBiRes *ReadBitResponce) Marshal() (payload []byte) {
	return
}

func (rBiRes *ReadBitResponce) Unmarshal(payload []byte) {
	rBiRes.bits = payload[8:]
}

func (rByRes *ReadByteResponce) Marshal() (payload []byte) {
	return
}

func (rByRes *ReadByteResponce) Unmarshal(payload []byte) {
	var payloadData [][]byte
	var workData []byte
	for currentIndex, currentBit := range payload[11:] {
		if currentIndex%2 == 0 {
			workData = []byte{currentBit}
		} else {
			workData = append(workData, currentBit)
			payloadData = append(payloadData, workData)
		}
	}
	rByRes.numberBits = payload[10]
	rByRes.data = payloadData
}
