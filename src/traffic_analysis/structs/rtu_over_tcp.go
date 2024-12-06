package structs

import (
	"fmt"
	"log"
	"modbus-emulator/conf"
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

func (eRes *RTUOverTCPErrorResponse) MarshalPayload() ([]uint16, error) {
	return []uint16{}, nil
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

func (req *RTUOverTCPRequest123456Response56) MarshalPayload() (payload []uint16, err error) {
	if req.HeaderError.FunctionID == conf.Functions.CoilsSimpleWrite {
		if req.ReadWriteDataHight == 255 {
			payload = append(payload, 1)
		} else {
			payload = append(payload, 0)
		}
	} else {
		var currentByte uint64
		if currentByte, err = strconv.ParseUint(fmt.Sprintf("%s%s",
			strconv.FormatUint(uint64(req.ReadWriteDataHight), 2), strconv.FormatUint(uint64(req.ReadWriteDataLow), 2)), 2, 64); err != nil {
			err = fmt.Errorf("error on marshaling registers data: %s", err)
			return
		}
		payload = append(payload, uint16(currentByte))
	}
	return
}

func (req *RTUOverTCPRequest123456Response56) LogPrint() {
	req.HeaderError.LogPrint()
	log.Printf("   Start address hight: %d", req.StartingAddressHight)
	log.Printf("   Start address low: %d", req.StartingAddressLow)
	if slices.Contains([]uint16{
		conf.Functions.CoilsRead,
		conf.Functions.DIRead,
		conf.Functions.HRRead,
		conf.Functions.IRRead}, req.HeaderError.FunctionID) {
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
	if slices.Contains([]uint16{conf.Functions.CoilsSimpleWrite, conf.Functions.HRSimpleWrite}, req.HeaderError.FunctionID) {
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

func (rRes *RTUOverTCPReadResponse) MarshalPayload() (payload []uint16, err error) {
	if slices.Contains([]uint16{conf.Functions.CoilsRead, conf.Functions.DIRead}, rRes.HeaderError.FunctionID) {
		if payload, err = InputsPayloadPreprocessing(rRes.Data); err != nil {
			err = fmt.Errorf("error on marshaling read data: %s", err)
		}
	} else {
		if payload, err = RegistersPayloadPreprocessing(rRes.Data); err != nil {
			err = fmt.Errorf("error on marshaling read data: %s", err)
		}
	}
	return
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

func (mWRes *RTUOverTCPMultipleWriteResponse) MarshalPayload() ([]uint16, error) {
	return []uint16{}, nil
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

func (mWReq *RTUOverTCPMultipleWriteRequest) MarshalPayload() (payload []uint16, err error) {
	if mWReq.Body.HeaderError.FunctionID == conf.Functions.CoilsMultipleWrite {
		countPayloadByte := int(mWReq.Body.QuantityOfRegistersHight + mWReq.Body.QuantityOfRegistersLow)
		for _, currentByte := range mWReq.Data {
			var currentBinaryData []string
			if countPayloadByte >= 8 {
				currentBinaryData = strings.Split(fmt.Sprintf("%08b", currentByte), "")
				countPayloadByte -= len(currentBinaryData)
			} else {
				currentBinaryData = strings.Split(strconv.FormatUint(uint64(currentByte), 2), "")
			}
			if len(currentBinaryData) != countPayloadByte {
				for {
					if len(currentBinaryData) == countPayloadByte {
						break
					}
					currentBinaryData = slices.Insert(currentBinaryData, 0, "0")
				}
			}
			for currentIndex := len(currentBinaryData) - 1; currentIndex > -1; currentIndex-- {
				var currentIntBuffer int
				if currentIntBuffer, err = strconv.Atoi(currentBinaryData[currentIndex]); err != nil {
					err = fmt.Errorf("error on marshaling coils write data: %s", err)
					return
				}
				payload = append(payload, uint16(currentIntBuffer))
			}
		}
		return
	}
	if payload, err = RegistersPayloadPreprocessing(mWReq.Data); err != nil {
		err = fmt.Errorf("error on marshaling HR write data: %s", err)
		return
	}
	return
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
