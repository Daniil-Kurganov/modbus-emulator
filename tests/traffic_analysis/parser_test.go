package trafficanalysis_test

import (
	"testing"
	"time"

	ta "modbus-emulator/src/traffic_analysis"

	"github.com/stretchr/testify/assert"
)

func TestParsePackets(t *testing.T) {
	testTable := []struct {
		typeObject      string
		filename        string
		expectedHistory []ta.History
	}{
		{
			typeObject: "",
			filename:   "",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-1",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 1},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  1,
							},
							AddressStart: []byte{0, 0},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 1},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  1,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits:       0,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 20, 974027915, time.Local),
				},
				{
					TransactionID: "0-2",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 2},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  2,
							},
							AddressStart: []byte{0, 16},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 2},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 2},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  2,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits:       0,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 21, 474872690, time.Local),
				},
				{
					TransactionID: "0-3",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 3},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  6,
							},
							AddressStart: []byte{0, 8},
							Data: &ta.WriteSimpleRequest{
								Payload: []byte{0, 39},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 3},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  6,
							},
							Data: &ta.WriteSimpleResponse{
								AddressStart: []byte{0, 8},
								WrittenBits:  []byte{0, 39},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 22, 377046677, time.Local),
				},
				{
					TransactionID: "0-4",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 4},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  1,
							},
							AddressStart: []byte{0, 5},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 5},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 4},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  1,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits:       0,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 23, 378031897, time.Local),
				},
				{
					TransactionID: "0-5",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 5},
								Protocol:      "modbus",
								BodyLength:    11,
								UnitID:        0,
								FunctionType:  15,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.WriteMultipleRequest{
								NumberRegisters: []byte{0, 4},
								NumberBits:      4,
								Data:            []byte{1, 1, 0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 5},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  15,
							},
							Data: &ta.WriteMultipleResponse{
								AddressStart:           []byte{0, 4},
								NumberWrittenRegisters: []byte{0, 4},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 23, 478357374, time.Local),
				},
			},
		},
	}
	var currentRecievedHistory []ta.History
	var err error
	for _, currentTestCase := range testTable {
		if currentRecievedHistory, err = ta.ParsePackets(); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		assert.Equalf(t, currentTestCase.expectedHistory, currentRecievedHistory,
			"Error: recieved and expected histories isn't equal:\n expected: %v;\n recieved: %v",
			currentTestCase.expectedHistory, currentRecievedHistory)
	}
}
