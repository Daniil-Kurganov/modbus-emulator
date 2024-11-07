package trafficanalysis_test

import (
	"io"
	"log"
	ta "modbus-emulator/src/traffic_analysis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestMBAPHeaderUnmarshal(t *testing.T) {
	testTable := []struct {
		payload        []byte
		expectedHeader ta.MBAPHeader
	}{
		{
			payload: []byte{0, 23, 0, 0, 0, 6, 0, 3, 0, 3, 0, 6},
			expectedHeader: ta.MBAPHeader{
				TransactionID: []byte{0, 23},
				Protocol:      "modbus",
				BodyLength:    6,
				UnitID:        0,
				FunctionType:  3,
			},
		},
		{
			payload:        []byte{},
			expectedHeader: ta.MBAPHeader{},
		},
	}
	for _, currentTestCase := range testTable {
		currentRecievedHeader := ta.MBAPHeader{}
		currentRecievedHeader.Unmarshal(currentTestCase.payload)
		assert.Equalf(t, currentTestCase.expectedHeader, currentRecievedHeader,
			"Error: expected and recieved headers isn't equal:\n expected: %v;\n recieved: %v", currentTestCase.expectedHeader, currentRecievedHeader)

	}
}

func TestTCPPacketRequestUnmarshalHeader(t *testing.T) {
	testTable := []struct {
		payload        []byte
		expectedPacket ta.TCPPacketRequest
	}{
		{
			payload: []byte{0, 23, 0, 0, 0, 6, 0, 3, 0, 3, 0, 6},
			expectedPacket: ta.TCPPacketRequest{
				Header: ta.MBAPHeader{
					TransactionID: []byte{0, 23},
					Protocol:      "modbus",
					BodyLength:    6,
					UnitID:        0,
					FunctionType:  3,
				},
				AddressStart: []byte{0, 3},
				Data:         nil,
			},
		},
		{
			payload: []byte{},
			expectedPacket: ta.TCPPacketRequest{
				Header: ta.MBAPHeader{},
				AddressStart: nil,
				Data: nil,
			},
		},
	}
	for _, currentTestCase := range testTable {
		currentRecievedPacket := ta.TCPPacketRequest{}
		currentRecievedPacket.UnmarshalHeader(currentTestCase.payload)
		assert.Equalf(t, currentTestCase.expectedPacket, currentRecievedPacket,
			"Error: expected and recieved packets isn't equal:\n expected: %v;\n recieved: %v", currentTestCase.expectedPacket, currentRecievedPacket)
	}
}

func TestTCPPacketResponseUnmarshalHeader(t *testing.T) {
	testTable := []struct{
		payload []byte
		expectedPacket ta.TCPPacketResponse
	}{
		{
			payload: []byte{0, 23, 0, 0, 0, 6, 0, 3, 0, 3, 0, 6},
			expectedPacket: ta.TCPPacketResponse{
				Header: ta.MBAPHeader{
					TransactionID: []byte{0, 23},
					Protocol:      "modbus",
					BodyLength:    6,
					UnitID:        0,
					FunctionType:  3,
				},
				Data: nil,
			},
		},
		{
			payload: []byte{},
			expectedPacket: ta.TCPPacketResponse{
				Header: ta.MBAPHeader{},
				Data: nil,
			},
		},
	}
	for _, currentTestCase := range testTable {
		currentRecievedPacket := ta.TCPPacketResponse{}
		currentRecievedPacket.UnmarshalHeader(currentTestCase.payload)
		assert.Equalf(t, currentTestCase.expectedPacket, currentRecievedPacket,
			"Error: expected and recieved packets isn't equal:\n expected: %v;\n recieved: %v", currentTestCase.expectedPacket, currentRecievedPacket)
	}
}

func TestReadRequestUnmarshal(t *testing.T) {
	testTable := []struct{
		payload []byte
		expectedRequest ta.ReadRequest
	}{
		{
			payload: []byte{0, 23, 0, 0, 0, 6, 0, 3, 0, 3, 0, 6},
			expectedRequest: ta.ReadRequest{
				NumberReadingBits: []byte{0, 6},
			},
		},
		{
			payload: []byte{},
			expectedRequest: ta.ReadRequest{},
		},
	}
	for _, currentTestCase := range testTable {
		currentRecievedRequest := ta.ReadRequest{}
		currentRecievedRequest.Unmarshal(currentTestCase.payload)
		assert.Equalf(t, currentTestCase.expectedRequest, currentRecievedRequest,
			"Error: expected and recieved requests isn't equal:\n expected: %v;\n recieved: %v", currentTestCase.expectedRequest, currentRecievedRequest)
	}
}

func TestWriteSimpleRequestUnmarshal(t *testing.T) {
	testTable := []struct{
		payload []byte
		expectedRequest ta.WriteSimpleRequest
	}{
		{
			payload: []byte{0, 17, 0, 0, 0, 6, 0, 6, 0, 3, 0, 7},
			expectedRequest: ta.WriteSimpleRequest{
				Payload: []byte{0, 7},
			},
		},
		{
			payload: []byte{},
			expectedRequest: ta.WriteSimpleRequest{},
		},
	}
	for _, currentTestCase := range testTable {
		currentRecievedRequest := ta.WriteSimpleRequest{}
		currentRecievedRequest.Unmarshal(currentTestCase.payload)
		assert.Equalf(t, currentTestCase.expectedRequest, currentRecievedRequest,
			"Error: expected and recieved requests isn't equal:\n expected: %v;\n recieved: %v", currentTestCase.expectedRequest, currentRecievedRequest)
	}
}

func TestWriteMultipleRequestUnmarshal(t *testing.T) {
	testTable := []struct{
		payload []byte
		expectedRequest ta.WriteMultipleRequest
	}{
		{
			payload: []byte{0, 2, 0, 0, 0, 11, 0, 16, 0, 3, 0, 2, 4, 0, 34, 0, 10},
			expectedRequest: ta.WriteMultipleRequest{
				NumberRegisters: []byte{0, 2},
				NumberBits: 4,
				Data: []byte{0, 34, 0, 10},
			},
		},
		{
			payload: []byte{},
			expectedRequest: ta.WriteMultipleRequest{},
		},
	}
	for _, currentTestCase := range testTable {
		currentRecievedRequest := ta.WriteMultipleRequest{}
		currentRecievedRequest.Unmarshal(currentTestCase.payload)
		assert.Equalf(t, currentTestCase.expectedRequest, currentRecievedRequest,
			"Error: expected and recieved requests isn't equal:\n expected: %v;\n recieved: %v", currentTestCase.expectedRequest, currentRecievedRequest)
	}
}
