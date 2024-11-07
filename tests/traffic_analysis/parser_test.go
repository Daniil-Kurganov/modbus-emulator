package trafficanalysis_test

import (
	"testing"

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
			typeObject: "coils",
			filename:   "read_41",
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
							AddressStart: []byte{0, 4},
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
								Bits:       1,
							},
						},
					},
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
								FunctionType:  1,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 2},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  1,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits:       1,
							},
						},
					},
				},
			},
		},
		{
			typeObject: "coils",
			filename:   "write_31",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-13",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  5,
							},
							AddressStart: []byte{0, 3},
							Data: &ta.WriteSimpleRequest{
								Payload: []byte{0, 0},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  5,
							},
							Data: &ta.WriteSimpleResponse{
								AddressStart: []byte{0, 3},
								WrittenBits:  []byte{0, 0},
							},
						},
					},
				},
				{
					TransactionID: "0-14",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  5,
							},
							AddressStart: []byte{0, 3},
							Data: &ta.WriteSimpleRequest{
								Payload: []byte{255, 0},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  5,
							},
							Data: &ta.WriteSimpleResponse{
								AddressStart: []byte{0, 3},
								WrittenBits:  []byte{255, 0},
							},
						},
					},
				},
			},
		},
		{
			typeObject: "coils",
			filename:   "write_34",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-15",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 15},
								Protocol:      "modbus",
								BodyLength:    11,
								UnitID:        0,
								FunctionType:  15,
							},
							AddressStart: []byte{0, 3},
							Data: &ta.WriteMultipleRequest{
								NumberRegisters: []byte{0, 4},
								NumberBits:      4,
								Data:            []byte{0, 1, 1, 0},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 15},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  15,
							},
							Data: &ta.WriteMultipleResponse{
								AddressStart:           []byte{0, 3},
								NumberWrittenRegisters: []byte{0, 4},
							},
						},
					},
				},
			},
		},
		{
			typeObject: "DI",
			filename:   "read_41",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-13",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  2,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  2,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits: 1,
							},
						},
					},
				},
				{
					TransactionID: "0-14",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  2,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  2,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits: 1,
							},
						},
					},
				},
				{
					TransactionID: "0-15",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 15},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  2,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 15},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  2,
							},
							Data: &ta.ReadBitResponse{
								NumberBits: 1,
								Bits: 1,
							},
						},
					},
				},
			},
		},
	}
	var currentRecievedHistory []ta.History
	var err error
	for _, currentTestCase := range testTable {
		if currentRecievedHistory, err = ta.ParsePackets("test_files", currentTestCase.typeObject, currentTestCase.filename); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		assert.Equalf(t, currentTestCase.expectedHistory, currentRecievedHistory,
			"Error: recieved and expected histories isn't equal:\n expected: %v;\n recieved: %v",
			currentTestCase.expectedHistory, currentRecievedHistory)
	}
}
