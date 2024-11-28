package tests_test

import (
	"testing"
	"time"

	ta "modbus-emulator/src/traffic_analysis"
	structs "modbus-emulator/src/traffic_analysis/structs"
	"modbus-emulator/src/utils"

	"github.com/stretchr/testify/assert"
)

func TestParsePackets(t *testing.T) {
	testTable := []struct {
		mode            string
		directoryPath   string
		expectedHistory []structs.HistoryEvent
	}{
		{
			mode:          "tcp",
			directoryPath: `src/pcapng_files/tests_files`,
			expectedHistory: []structs.HistoryEvent{
				{
					TransactionID: "0-1",
					Handshake: structs.Handshake{
						Request: &structs.TCPRequest{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 1},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  1,
							},
							AddressStart: []byte{0, 0},
							Data: &structs.TCPReadRequest{
								NumberReadingBits: []byte{0, 1},
							},
						},
						Response: &structs.TCPResponse{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 1},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  1,
							},
							Data: &structs.TCPReadBitResponse{
								NumberBits: 1,
								Bits:       0,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 20, 974027915, time.Local),
				},
				{
					TransactionID: "0-2",
					Handshake: structs.Handshake{
						Request: &structs.TCPRequest{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 2},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  2,
							},
							AddressStart: []byte{0, 16},
							Data: &structs.TCPReadRequest{
								NumberReadingBits: []byte{0, 2},
							},
						},
						Response: &structs.TCPResponse{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 2},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  2,
							},
							Data: &structs.TCPReadBitResponse{
								NumberBits: 1,
								Bits:       0,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 21, 474872690, time.Local),
				},
				{
					TransactionID: "0-3",
					Handshake: structs.Handshake{
						Request: &structs.TCPRequest{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 3},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  6,
							},
							AddressStart: []byte{0, 8},
							Data: &structs.TCPWriteSimpleRequest{
								Payload: []byte{0, 39},
							},
						},
						Response: &structs.TCPResponse{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 3},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  6,
							},
							Data: &structs.TCPWriteSimpleResponse{
								AddressStart: []byte{0, 8},
								WrittenBits:  []byte{0, 39},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 22, 377046677, time.Local),
				},
				{
					TransactionID: "0-4",
					Handshake: structs.Handshake{
						Request: &structs.TCPRequest{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 4},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  1,
							},
							AddressStart: []byte{0, 5},
							Data: &structs.TCPReadRequest{
								NumberReadingBits: []byte{0, 5},
							},
						},
						Response: &structs.TCPResponse{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 4},
								Protocol:      "modbus",
								BodyLength:    4,
								UnitID:        0,
								FunctionType:  1,
							},
							Data: &structs.TCPReadBitResponse{
								NumberBits: 1,
								Bits:       0,
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 23, 378031897, time.Local),
				},
				{
					TransactionID: "0-5",
					Handshake: structs.Handshake{
						Request: &structs.TCPRequest{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 5},
								Protocol:      "modbus",
								BodyLength:    11,
								UnitID:        0,
								FunctionType:  15,
							},
							AddressStart: []byte{0, 4},
							Data: &structs.TCPWriteMultipleRequest{
								NumberRegisters: []byte{0, 4},
								NumberBits:      4,
								Data:            []byte{1, 1, 0, 1},
							},
						},
						Response: &structs.TCPResponse{
							Header: structs.MBAPHeader{
								TransactionID: []byte{0, 5},
								Protocol:      "modbus",
								BodyLength:    6,
								UnitID:        0,
								FunctionType:  15,
							},
							Data: &structs.TCPWriteMultipleResponse{
								AddressStart:           []byte{0, 4},
								NumberWrittenRegisters: []byte{0, 4},
							},
						},
					},
					TransactionTime: time.Date(2024, 11, 11, 12, 53, 23, 478357374, time.Local),
				},
			},
		},
		{
			mode:          "rtu_over_tcp",
			directoryPath: `src/pcapng_files/tests_files`,
			expectedHistory: []structs.HistoryEvent{
				{
					TransactionID: "1",
					Handshake: structs.Handshake{
						Request: &structs.RTUOverTCPRequest123456Response56{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      1,
								ErrorCheckLow:   253,
								ErrorCheckHight: 202,
							},
							StartingAddressHight: 0,
							StartingAddressLow:   0,
							ReadWriteDataHight:   0,
							ReadWriteDataLow:     1,
						},
						Response: &structs.RTUOverTCPReadResponse{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      1,
								ErrorCheckLow:   81,
								ErrorCheckHight: 136,
							},
							ByteCount: 1,
							Data:      []uint16{0},
						},
					},
					TransactionTime: time.Date(2024, 11, 20, 12, 31, 18, 925704421, time.Local),
				},
				{
					TransactionID: "2",
					Handshake: structs.Handshake{
						Request: &structs.RTUOverTCPRequest123456Response56{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      2,
								ErrorCheckLow:   248,
								ErrorCheckHight: 14,
							},
							StartingAddressHight: 0,
							StartingAddressLow:   16,
							ReadWriteDataHight:   0,
							ReadWriteDataLow:     2,
						},
						Response: &structs.RTUOverTCPReadResponse{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      2,
								ErrorCheckLow:   161,
								ErrorCheckHight: 136,
							},
							ByteCount: 1,
							Data:      []uint16{0},
						},
					},
					TransactionTime: time.Date(2024, 11, 20, 12, 31, 19, 432394852, time.Local),
				},
				{
					TransactionID: "3",
					Handshake: structs.Handshake{
						Request: &structs.RTUOverTCPRequest123456Response56{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      6,
								ErrorCheckLow:   72,
								ErrorCheckHight: 18,
							},
							StartingAddressHight: 0,
							StartingAddressLow:   8,
							ReadWriteDataHight:   0,
							ReadWriteDataLow:     39,
						},
						Response: &structs.RTUOverTCPRequest123456Response56{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      6,
								ErrorCheckLow:   72,
								ErrorCheckHight: 18,
							},
							StartingAddressHight: 0,
							StartingAddressLow:   8,
							ReadWriteDataHight:   0,
							ReadWriteDataLow:     39,
						},
					},
					TransactionTime: time.Date(2024, 11, 20, 12, 31, 20, 341282997, time.Local),
				},
				{
					TransactionID: "4",
					Handshake: structs.Handshake{
						Request: &structs.RTUOverTCPRequest123456Response56{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      1,
								ErrorCheckLow:   236,
								ErrorCheckHight: 8,
							},
							StartingAddressHight: 0,
							StartingAddressLow:   5,
							ReadWriteDataHight:   0,
							ReadWriteDataLow:     5,
						},
						Response: &structs.RTUOverTCPReadResponse{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      1,
								ErrorCheckLow:   145,
								ErrorCheckHight: 139,
							},
							ByteCount: 1,
							Data:      []uint16{5},
						},
					},
					TransactionTime: time.Date(2024, 11, 20, 12, 31, 21, 349492607, time.Local),
				},
				{
					TransactionID: "5",
					Handshake: structs.Handshake{
						Request: &structs.RTUOverTCPMultipleWriteRequest{
							Body: structs.RTUOverTCPMultipleWriteResponse{
								HeaderError: structs.HeaderErrorCheck{
									SlaveAddress:    1,
									FunctionID:      15,
									ErrorCheckLow:   142,
									ErrorCheckHight: 145,
								},
								RegisterAddressHight:     0,
								RegisterAddressLow:       4,
								QuantityOfRegistersHight: 0,
								QuantityOfRegistersLow:   4,
							},
							ByteCount: 1,
							Data:      []uint16{11},
						},
						Response: &structs.RTUOverTCPMultipleWriteResponse{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      15,
								ErrorCheckLow:   21,
								ErrorCheckHight: 201,
							},
							RegisterAddressHight:     0,
							RegisterAddressLow:       4,
							QuantityOfRegistersHight: 0,
							QuantityOfRegistersLow:   4,
						},
					},
					TransactionTime: time.Date(2024, 11, 20, 12, 31, 21, 456314271, time.Local),
				},
				{
					TransactionID: "6",
					Handshake: structs.Handshake{
						Request: &structs.RTUOverTCPRequest123456Response56{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      3,
								ErrorCheckLow:   180,
								ErrorCheckHight: 9,
							},
							StartingAddressHight: 0,
							StartingAddressLow:   3,
							ReadWriteDataHight:   0,
							ReadWriteDataLow:     4,
						},
						Response: &structs.RTUOverTCPReadResponse{
							HeaderError: structs.HeaderErrorCheck{
								SlaveAddress:    1,
								FunctionID:      3,
								ErrorCheckLow:   84,
								ErrorCheckHight: 20,
							},
							ByteCount: 8,
							Data:      []uint16{0, 0, 0, 16, 0, 0, 0, 0},
						},
					},
					TransactionTime: time.Date(2024, 11, 20, 12, 31, 22, 485545105, time.Local),
				},
			},
		},
	}
	var currentRecievedHistory []structs.HistoryEvent
	var err error
	for _, currentTestCase := range testTable {
		utils.WorkMode = currentTestCase.mode
		utils.DumpDirectoryPath = currentTestCase.directoryPath
		if currentRecievedHistory, err = ta.ParseDump(); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		assert.Equalf(t, currentTestCase.expectedHistory, currentRecievedHistory,
			"Error: recieved and expected histories isn't equal:\n expected: %v;\n recieved: %v",
			currentTestCase.expectedHistory, currentRecievedHistory)
	}
}
