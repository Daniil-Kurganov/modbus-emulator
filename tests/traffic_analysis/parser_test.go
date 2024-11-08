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
					TransactionTime: time.Date(2024, 11, 7, 16, 52, 11, 968278907, time.Local),
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
					TransactionTime: time.Date(2024, 11, 7, 16, 52, 12, 469062685, time.Local),
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
					TransactionTime: time.Date(2024, 11, 7, 16, 58, 6, 734201161, time.Local),
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
					TransactionTime: time.Date(2024, 11, 7, 16, 58, 7, 234978658, time.Local),
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
					TransactionTime: time.Date(2024, 11, 7, 16, 59, 9, 885417797, time.Local),
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
								Bits:       1,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 0, 34, 707271279, time.Local),
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
								Bits:       1,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 0, 35, 207995000, time.Local),
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
								Bits:       1,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 0, 35, 708765041, time.Local),
				},
			},
		},
		{
			typeObject: "HR",
			filename:   "read_41",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-8",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 8},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  3,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 8},
								Protocol:      "modbus",
								BodyLength:    5,
								UnitID:        0,
								FunctionType:  3,
							},
							Data: &ta.ReadByteResponse{
								NumberBits: 2,
								Data:       [][]byte{{0, 6}},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 2, 9, 648349087, time.Local),
				},
				{
					TransactionID: "0-9",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 9},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  3,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 9},
								Protocol:      "modbus",
								BodyLength:    5,
								UnitID:        0,
								FunctionType:  3,
							},
							Data: &ta.ReadByteResponse{
								NumberBits: 2,
								Data:       [][]byte{{0, 8}},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 2, 10, 149119047, time.Local),
				},
			},
		},
		{
			typeObject: "HR",
			filename:   "write_41",
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
								FunctionType:  6,
							},
							AddressStart: []byte{0, 4},
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
								FunctionType:  6,
							},
							Data: &ta.WriteSimpleResponse{
								AddressStart: []byte{0, 4},
								WrittenBits:  []byte{0, 0},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 3, 29, 938303882, time.Local),
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
								FunctionType:  6,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.WriteSimpleRequest{
								Payload: []byte{0, 11},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  6,
							},
							Data: &ta.WriteSimpleResponse{
								AddressStart: []byte{0, 4},
								WrittenBits:  []byte{0, 11},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 3, 30, 439083007, time.Local),
				},
			},
		},
		{
			typeObject: "HR",
			filename:   "write_42",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-13",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    11,
								UnitID:        0,
								FunctionType:  16,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.WriteMultipleRequest{
								NumberRegisters: []byte{0, 2},
								NumberBits:      4,
								Data:            []byte{0, 11, 0, 20},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  16,
							},
							Data: &ta.WriteMultipleResponse{
								AddressStart:           []byte{0, 4},
								NumberWrittenRegisters: []byte{0, 2},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 4, 34, 332275682, time.Local),
				},
				{
					TransactionID: "0-14",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    11,
								UnitID:        0,
								FunctionType:  16,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.WriteMultipleRequest{
								NumberRegisters: []byte{0, 2},
								NumberBits:      4,
								Data:            []byte{0, 11, 0, 20},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 14},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  16,
							},
							Data: &ta.WriteMultipleResponse{
								AddressStart:           []byte{0, 4},
								NumberWrittenRegisters: []byte{0, 2},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 4, 34, 833100830, time.Local),
				},
			},
		},
		{
			typeObject: "IR",
			filename:   "read_41",
			expectedHistory: []ta.History{
				{
					TransactionID: "0-11",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 11},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  4,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 11},
								Protocol:      "modbus",
								BodyLength:    5,
								UnitID:        0,
								FunctionType:  4,
							},
							Data: &ta.ReadByteResponse{
								NumberBits: 2,
								Data:       [][]byte{{0, 6}},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 5, 37, 350750846, time.Local),
				},
				{
					TransactionID: "0-12",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 12},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  4,
							},
							AddressStart: []byte{0, 4},
							Data: &ta.ReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &ta.TCPPacketResponse{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 12},
								Protocol:      "modbus",
								BodyLength:    5,
								UnitID:        0,
								FunctionType:  4,
							},
							Data: &ta.ReadByteResponse{
								NumberBits: 2,
								Data:       [][]byte{{0, 6}},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 5, 37, 851532034, time.Local),
				},
				{
					TransactionID: "0-13",
					Handshake: ta.Handshake{
						Request: &ta.TCPPacketRequest{
							Header: ta.MBAPHeader{
								TransactionID: []byte{0, 13},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  4,
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
								BodyLength:    5,
								UnitID:        0,
								FunctionType:  4,
							},
							Data: &ta.ReadByteResponse{
								NumberBits: 2,
								Data:       [][]byte{{0, 6}},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 7, 17, 5, 38, 352389461, time.Local),
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
