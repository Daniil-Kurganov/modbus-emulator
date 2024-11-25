package structs

import (
	"log"
	"slices"
	"strconv"
	"strings"
)

type (
	HeaderErrorCheck struct {
		SlaveAddress    uint16
		FunctionID      uint16
		ErrorCheckLow   uint16
		ErrorCheckHight uint16
	}
	RTUOverTCPErrorResponse struct {
		HeaderError HeaderErrorCheck
		ErrorCode   uint16
	}
	RTUOverTCPRequest123456Response56 struct {
		HeaderError          HeaderErrorCheck
		StartingAddressHight uint16
		StartingAddressLow   uint16
		ReadWriteDataHight   uint16
		ReadWriteDataLow     uint16
	}
	RTUOverTCPReadResponse struct {
		HeaderError HeaderErrorCheck
		ByteCount   uint16
		Data        []uint16
	}
	RTUOverTCPMultipleWriteRequest struct {
		Body      RTUOverTCPMultipleWriteResponse
		ByteCount uint16
		Data      []uint16
	}
	RTUOverTCPMultipleWriteResponse struct {
		HeaderError              HeaderErrorCheck
		RegisterAddressHight     uint16
		RegisterAddressLow       uint16
		QuantityOfRegistersHight uint16
		QuantityOfRegistersLow   uint16
	}
)

func (h *HeaderErrorCheck) Unmarshal(payload []byte) {
	h.SlaveAddress = uint16(payload[0])
	h.FunctionID = uint16(payload[1])
	h.ErrorCheckLow = uint16(payload[len(payload)-2])
	h.ErrorCheckHight = uint16(payload[len(payload)-1])
}

func (h *HeaderErrorCheck) LogPrint() {
	log.Printf("   Slave address: %d", h.SlaveAddress)
	log.Printf("   Function ID: %d", h.FunctionID)
	log.Printf("   Error check low: %d", h.ErrorCheckLow)
	log.Printf("   Error check hight: %d", h.ErrorCheckHight)
}

func (eRes *RTUOverTCPErrorResponse) Unmarshal(payload []byte) {
	eRes.HeaderError.Unmarshal(payload)
	eRes.ErrorCode = uint16(payload[2])
}

func (eRes *RTUOverTCPErrorResponse) MarshalPayload() []uint16 {
	return []uint16{}
}

func (eRes *RTUOverTCPErrorResponse) LogPrint() {
	eRes.HeaderError.LogPrint()
	log.Printf("   Error code: %d", eRes.ErrorCode)
}

func (eRes *RTUOverTCPErrorResponse) GetFunctionID() uint16 {
	return eRes.HeaderError.FunctionID
}

func (req *RTUOverTCPRequest123456Response56) Unmarshal(payload []byte) {
	req.HeaderError.Unmarshal(payload)
	req.StartingAddressHight = uint16(payload[2])
	req.StartingAddressLow = uint16(payload[3])
	req.ReadWriteDataHight = uint16(payload[4])
	req.ReadWriteDataLow = uint16(payload[5])
}

func (req *RTUOverTCPRequest123456Response56) MarshalPayload() []uint16 {
	if slices.Contains([]uint16{5, 6}, req.HeaderError.FunctionID) {
		return []uint16{req.ReadWriteDataHight, req.ReadWriteDataLow}
	}
	return []uint16{}
}

func (req *RTUOverTCPRequest123456Response56) LogPrint() {
	req.HeaderError.LogPrint()
	log.Printf("   Start address hight: %d", req.StartingAddressHight)
	log.Printf("   Start address low: %d", req.StartingAddressLow)
	if slices.Contains([]uint16{1, 2, 3, 4}, req.HeaderError.FunctionID) {
		log.Printf("   Quantity of registers hight: %d", req.ReadWriteDataHight)
		log.Printf("   Quantity of registers low: %d", req.ReadWriteDataLow)
	} else {
		log.Printf("   Write data hight: %d", req.ReadWriteDataHight)
		log.Printf("   Write data low: %d", req.ReadWriteDataLow)
	}
}

func (req *RTUOverTCPRequest123456Response56) MarshalAddress() []uint16 {
	return []uint16{req.StartingAddressHight, req.StartingAddressLow}
}

func (req *RTUOverTCPRequest123456Response56) MarshalQuantity() []uint16 {
	if slices.Contains([]uint16{5, 6}, req.HeaderError.FunctionID) {
		return []uint16{0, 1}
	}
	return []uint16{req.ReadWriteDataHight, req.ReadWriteDataLow}
}

func (req *RTUOverTCPRequest123456Response56) GetFunctionID() uint16 {
	return req.HeaderError.FunctionID
}

func (rRes *RTUOverTCPReadResponse) Unmarshal(payload []byte) {
	rRes.HeaderError.Unmarshal(payload)
	rRes.ByteCount = uint16(payload[2])
	for currentBitIndex := 3; currentBitIndex < 3+int(rRes.ByteCount); currentBitIndex++ {
		rRes.Data = append(rRes.Data, uint16(payload[currentBitIndex]))
	}
}

func (rRes *RTUOverTCPReadResponse) MarshalPayload() []uint16 {
	return []uint16{}
}

func (rRes *RTUOverTCPReadResponse) LogPrint() {
	rRes.HeaderError.LogPrint()
	log.Printf("   Byte count: %d", rRes.ByteCount)
	log.Printf("   Data: %v", rRes.Data)
}

func (rRes *RTUOverTCPReadResponse) GetFunctionID() uint16 {
	return rRes.HeaderError.FunctionID
}

func (mWRes *RTUOverTCPMultipleWriteResponse) Unmarshal(payload []byte) {
	mWRes.HeaderError.Unmarshal(payload)
	mWRes.HeaderError.Unmarshal(payload)
	mWRes.RegisterAddressHight = uint16(payload[2])
	mWRes.RegisterAddressLow = uint16(payload[3])
	mWRes.QuantityOfRegistersHight = uint16(payload[4])
	mWRes.QuantityOfRegistersLow = uint16(payload[5])
}

func (mWRes *RTUOverTCPMultipleWriteResponse) MarshalPayload() []uint16 {
	return []uint16{}
}

func (mWRes *RTUOverTCPMultipleWriteResponse) LogPrint() {
	mWRes.HeaderError.LogPrint()
	log.Printf("   Register address hight: %d", mWRes.RegisterAddressHight)
	log.Printf("   Register address low: %d", mWRes.RegisterAddressLow)
	log.Printf("   Quantity of registers hight: %d", mWRes.QuantityOfRegistersHight)
	log.Printf("   Quantity of registers low: %d", mWRes.QuantityOfRegistersLow)
}

func (mWRes *RTUOverTCPMultipleWriteResponse) GetFunctionID() uint16 {
	return mWRes.HeaderError.FunctionID
}

func (mWReq *RTUOverTCPMultipleWriteRequest) Unmarshal(payload []byte) {
	mWReq.Body.Unmarshal(payload)
	mWReq.ByteCount = uint16(payload[6])
	for currentBitIndex := 7; currentBitIndex < 7+int(mWReq.ByteCount); currentBitIndex++ {
		mWReq.Data = append(mWReq.Data, uint16(payload[currentBitIndex]))
	}
}

func (mWReq *RTUOverTCPMultipleWriteRequest) MarshalPayload() []uint16 {
	var err error
	if mWReq.Body.HeaderError.FunctionID == 15 {
		var payload []uint16
		for _, currentByte := range mWReq.Data {
			currentBinaryData := strings.Split(strconv.FormatUint(uint64(currentByte), 2), "")
			for currentIndex := len(currentBinaryData) - 1; currentIndex > -1; currentIndex-- {
				var currentIntBuffer int
				if currentIntBuffer, err = strconv.Atoi(currentBinaryData[currentIndex]); err != nil {
					log.Fatalf("Error on marshaling payload: %s", err)
				}
				payload = append(payload, uint16(currentIntBuffer))
			}
		}
		return payload
	}
	return []uint16{}
}

func (mWReq *RTUOverTCPMultipleWriteRequest) LogPrint() {
	mWReq.Body.LogPrint()
	log.Printf("   Byte count: %d", mWReq.ByteCount)
	log.Printf("   Data: %v", mWReq.Data)
}

func (mWReq *RTUOverTCPMultipleWriteRequest) MarshalAddress() []uint16 {
	return []uint16{mWReq.Body.RegisterAddressHight, mWReq.Body.RegisterAddressLow}
}

func (mWReq *RTUOverTCPMultipleWriteRequest) MarshalQuantity() []uint16 {
	return []uint16{mWReq.Body.QuantityOfRegistersHight, mWReq.Body.QuantityOfRegistersLow}
}
