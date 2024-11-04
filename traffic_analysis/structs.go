package trafficanalysis

type (
	Handshake struct {
		request  TCPPacket
		responce TCPPacket
	}
	TCPPacket struct {
		PacketNumber byte
		Protocol     string
		BodyLength   byte
		UnitID       byte
		FunctionType byte
		AddressStart []byte
		DataPayload  dataPayload
	}
	dataPayload interface {
		getOperationCode() int // 0 - read, 1 - write simple, 2 - write multiple
	}
	readRequest struct {
		operationCode int
		payload       []byte // number of reading bits
	}
	readResponce struct {
		operationCode int
		payload       struct {
			numberBits byte
			data       [][]byte // like: [[0, 26], [0, 130]]; len = numberBytes
		}
	}
	writeSimpleRequest struct {
		operationCode int
		payload       []byte // like: [0, 6]
	}
	writeMultipleRequest struct {
		operationCode int
		payload       struct {
			numberRegisters []byte // 2 bits, like: [0, 3] or [0, 2]
			numberBits      byte
			data            [][]byte // like: [[0, 45], [0, 35]]; len = numberRegisters[1]
		}
	}
	writeSingleResponce struct {
		operationCode int
		payload       []byte // written bits
	}
	writeMultipleResponce struct {
		operationCode int
		payload       byte // number of written bits
	}
)

func (rReq *readRequest) getOperationCode() int {
	return rReq.operationCode
}

func (rRes *readResponce) getOperationCode() int {
	return rRes.operationCode
}

func (wSReq *writeSimpleRequest) getOperationCode() int {
	return wSReq.operationCode
}

func (wMReq *writeMultipleRequest) getOperationCode() int {
	return wMReq.operationCode
}

func (wSRes *writeSingleResponce) getOperationCode() int {
	return wSRes.operationCode
}

func (wMRes *writeMultipleResponce) getOperationCode() int {
	return wMRes.operationCode
}
