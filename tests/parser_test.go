package tests_test

import (
	"testing"
	"time"

	"modbus-emulator/conf"
	ta "modbus-emulator/src/traffic_analysis"
	structs "modbus-emulator/src/traffic_analysis/structs"

	"github.com/stretchr/testify/assert"
)

func TestParsePackets(t *testing.T) {
	testTable := []struct {
		mode            string
		directoryPath   string
		ports           map[uint16]conf.ServerSocket
		expectedHistory map[uint16]structs.ServerHistory
	}{
		{
			mode:          "tcp",
			directoryPath: `pcapng_files/tests_files/simple_port`,
			ports: map[uint16]conf.ServerSocket{
				1502: {
					HostAddress: "127.0.0.1",
					PortAddress: "1502",
				},
			},
			expectedHistory: map[uint16]structs.ServerHistory{
				1502: {
					Transactions: []structs.HistoryEvent{
						{
							Header: structs.SlaveTransaction{
								SlaveID:       0,
								TransactionID: "0-1",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       0,
								TransactionID: "0-2",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       0,
								TransactionID: "0-3",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       0,
								TransactionID: "0-4",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       0,
								TransactionID: "0-5",
							},
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
					Slaves: []uint8{0},
				},
			},
		},
		{
			mode:          "rtu_over_tcp",
			directoryPath: `pcapng_files/tests_files/simple_port`,
			ports: map[uint16]conf.ServerSocket{
				1502: {
					HostAddress: "127.0.0.1",
					PortAddress: "1502",
				},
			},
			expectedHistory: map[uint16]structs.ServerHistory{
				1502: {
					Transactions: []structs.HistoryEvent{
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "1",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "2",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "3",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "4",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "5",
							},
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
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "6",
							},
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
					Slaves: []uint8{1},
				},
			},
		},
		{
			mode:          "rtu_over_tcp",
			directoryPath: `pcapng_files/tests_files/multiple_ports`,
			ports: map[uint16]conf.ServerSocket{
				1502: {
					HostAddress: "127.0.0.1",
					PortAddress: "1502",
				},
				1503: {
					HostAddress: "127.0.0.1",
					PortAddress: "1503",
				},
			},
			expectedHistory: map[uint16]structs.ServerHistory{
				1502: {
					Transactions: []structs.HistoryEvent{
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "1",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    1,
											FunctionID:      15,
											ErrorCheckLow:   122,
											ErrorCheckHight: 150,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       7,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   3,
									},
									ByteCount: 1,
									Data:      []uint16{3},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      15,
										ErrorCheckLow:   164,
										ErrorCheckHight: 11,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       7,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   3,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 451189674, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "2",
							},
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
										ErrorCheckLow:   81,
										ErrorCheckHight: 141,
									},
									ByteCount: 1,
									Data:      []uint16{12},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 461782916, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "3",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      2,
										ErrorCheckLow:   233,
										ErrorCheckHight: 207,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   9,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     11,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      2,
										ErrorCheckLow:   212,
										ErrorCheckHight: 216,
									},
									ByteCount: 2,
									Data:      []uint16{146, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 471331375, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "4",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    1,
											FunctionID:      16,
											ErrorCheckLow:   52,
											ErrorCheckHight: 33,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       150,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   7,
									},
									ByteCount: 14,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      16,
										ErrorCheckLow:   97,
										ErrorCheckHight: 231,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       150,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   7,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 481011269, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "5",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      3,
										ErrorCheckLow:   229,
										ErrorCheckHight: 226,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   150,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     15,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      3,
										ErrorCheckLow:   243,
										ErrorCheckHight: 180,
									},
									ByteCount: 30,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 498718091, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "6",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      4,
										ErrorCheckLow:   49,
										ErrorCheckHight: 198,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   4,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     18,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      4,
										ErrorCheckLow:   49,
										ErrorCheckHight: 236,
									},
									ByteCount: 36,
									Data:      []uint16{0, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 25, 248},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 508274173, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "7",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    2,
											FunctionID:      15,
											ErrorCheckLow:   58,
											ErrorCheckHight: 131,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       7,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   3,
									},
									ByteCount: 1,
									Data:      []uint16{3},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      15,
										ErrorCheckLow:   164,
										ErrorCheckHight: 56,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       7,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   3,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 517965863, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "8",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      1,
										ErrorCheckLow:   236,
										ErrorCheckHight: 59,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   5,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     5,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      1,
										ErrorCheckLow:   81,
										ErrorCheckHight: 201,
									},
									ByteCount: 1,
									Data:      []uint16{12},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 528515109, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "9",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      2,
										ErrorCheckLow:   233,
										ErrorCheckHight: 252,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   9,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     11,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      2,
										ErrorCheckLow:   253,
										ErrorCheckHight: 184,
									},
									ByteCount: 2,
									Data:      []uint16{0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 538020401, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "10",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    2,
											FunctionID:      16,
											ErrorCheckLow:   7,
											ErrorCheckHight: 18,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       150,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   7,
									},
									ByteCount: 14,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      16,
										ErrorCheckLow:   97,
										ErrorCheckHight: 212,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       150,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   7,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 547629103, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "11",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      3,
										ErrorCheckLow:   229,
										ErrorCheckHight: 209,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   150,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     15,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      3,
										ErrorCheckLow:   71,
										ErrorCheckHight: 180,
									},
									ByteCount: 30,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 565165744, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "12",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      4,
										ErrorCheckLow:   49,
										ErrorCheckHight: 245,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   4,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     18,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      4,
										ErrorCheckLow:   145,
										ErrorCheckHight: 233,
									},
									ByteCount: 36,
									Data:      []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 574719199, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "13",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    3,
											FunctionID:      15,
											ErrorCheckLow:   251,
											ErrorCheckHight: 79,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       7,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   3,
									},
									ByteCount: 1,
									Data:      []uint16{3},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      15,
										ErrorCheckLow:   165,
										ErrorCheckHight: 233,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       7,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   3,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 584270297, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "14",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      1,
										ErrorCheckLow:   237,
										ErrorCheckHight: 234,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   5,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     5,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      1,
										ErrorCheckLow:   80,
										ErrorCheckHight: 53,
									},
									ByteCount: 1,
									Data:      []uint16{12},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 594859677, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "15",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      2,
										ErrorCheckLow:   232,
										ErrorCheckHight: 45,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   9,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     11,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      2,
										ErrorCheckLow:   192,
										ErrorCheckHight: 120,
									},
									ByteCount: 2,
									Data:      []uint16{0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 604409306, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "16",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    3,
											FunctionID:      16,
											ErrorCheckLow:   23,
											ErrorCheckHight: 195,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       150,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   7,
									},
									ByteCount: 14,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      16,
										ErrorCheckLow:   96,
										ErrorCheckHight: 5,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       150,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   7,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 613889633, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "17",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      3,
										ErrorCheckLow:   228,
										ErrorCheckHight: 0,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   150,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     15,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      3,
										ErrorCheckLow:   42,
										ErrorCheckHight: 116,
									},
									ByteCount: 30,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 631491015, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "18",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      4,
										ErrorCheckLow:   48,
										ErrorCheckHight: 36,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   4,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     18,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      4,
										ErrorCheckLow:   153,
										ErrorCheckHight: 69,
									},
									ByteCount: 36,
									Data:      []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 641006079, time.Local),
						},
					},
					Slaves: []uint8{1, 2, 3},
				},
				1503: {
					Transactions: []structs.HistoryEvent{
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "1",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    1,
											FunctionID:      15,
											ErrorCheckLow:   122,
											ErrorCheckHight: 150,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       7,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   3,
									},
									ByteCount: 1,
									Data:      []uint16{3},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      15,
										ErrorCheckLow:   164,
										ErrorCheckHight: 11,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       7,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   3,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 451153346, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "2",
							},
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
										ErrorCheckLow:   81,
										ErrorCheckHight: 141,
									},
									ByteCount: 1,
									Data:      []uint16{12},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 461764762, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "3",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      2,
										ErrorCheckLow:   233,
										ErrorCheckHight: 207,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   9,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     11,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      2,
										ErrorCheckLow:   212,
										ErrorCheckHight: 216,
									},
									ByteCount: 2,
									Data:      []uint16{146, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 471342633, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "4",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    1,
											FunctionID:      16,
											ErrorCheckLow:   52,
											ErrorCheckHight: 33,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       150,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   7,
									},
									ByteCount: 14,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      16,
										ErrorCheckLow:   97,
										ErrorCheckHight: 231,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       150,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   7,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 481123410, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "5",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      3,
										ErrorCheckLow:   229,
										ErrorCheckHight: 226,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   150,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     15,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      3,
										ErrorCheckLow:   243,
										ErrorCheckHight: 180,
									},
									ByteCount: 30,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 498709665, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "6",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      4,
										ErrorCheckLow:   49,
										ErrorCheckHight: 198,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   4,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     18,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    1,
										FunctionID:      4,
										ErrorCheckLow:   49,
										ErrorCheckHight: 236,
									},
									ByteCount: 36,
									Data:      []uint16{0, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 25, 248},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 508274805, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "7",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    2,
											FunctionID:      15,
											ErrorCheckLow:   58,
											ErrorCheckHight: 131,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       7,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   3,
									},
									ByteCount: 1,
									Data:      []uint16{3},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      15,
										ErrorCheckLow:   164,
										ErrorCheckHight: 56,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       7,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   3,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 518013379, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "8",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      1,
										ErrorCheckLow:   236,
										ErrorCheckHight: 59,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   5,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     5,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      1,
										ErrorCheckLow:   81,
										ErrorCheckHight: 201,
									},
									ByteCount: 1,
									Data:      []uint16{12},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 528485290, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "9",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      2,
										ErrorCheckLow:   233,
										ErrorCheckHight: 252,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   9,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     11,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      2,
										ErrorCheckLow:   253,
										ErrorCheckHight: 184,
									},
									ByteCount: 2,
									Data:      []uint16{0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 537991649, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "10",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    2,
											FunctionID:      16,
											ErrorCheckLow:   7,
											ErrorCheckHight: 18,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       150,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   7,
									},
									ByteCount: 14,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      16,
										ErrorCheckLow:   97,
										ErrorCheckHight: 212,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       150,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   7,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 547639948, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "11",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      3,
										ErrorCheckLow:   229,
										ErrorCheckHight: 209,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   150,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     15,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      3,
										ErrorCheckLow:   71,
										ErrorCheckHight: 180,
									},
									ByteCount: 30,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 565120052, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "12",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      4,
										ErrorCheckLow:   49,
										ErrorCheckHight: 245,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   4,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     18,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    2,
										FunctionID:      4,
										ErrorCheckLow:   145,
										ErrorCheckHight: 233,
									},
									ByteCount: 36,
									Data:      []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 574742165, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "13",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    3,
											FunctionID:      15,
											ErrorCheckLow:   251,
											ErrorCheckHight: 79,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       7,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   3,
									},
									ByteCount: 1,
									Data:      []uint16{3},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      15,
										ErrorCheckLow:   165,
										ErrorCheckHight: 233,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       7,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   3,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 584320110, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "14",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      1,
										ErrorCheckLow:   237,
										ErrorCheckHight: 234,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   5,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     5,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      1,
										ErrorCheckLow:   80,
										ErrorCheckHight: 53,
									},
									ByteCount: 1,
									Data:      []uint16{12},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 594861107, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "15",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      2,
										ErrorCheckLow:   232,
										ErrorCheckHight: 45,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   9,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     11,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      2,
										ErrorCheckLow:   192,
										ErrorCheckHight: 120,
									},
									ByteCount: 2,
									Data:      []uint16{0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 604416840, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "16",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPMultipleWriteRequest{
									Body: structs.RTUOverTCPMultipleWriteResponse{
										HeaderError: structs.HeaderErrorCheck{
											SlaveAddress:    3,
											FunctionID:      16,
											ErrorCheckLow:   23,
											ErrorCheckHight: 195,
										},
										RegisterAddressHight:     0,
										RegisterAddressLow:       150,
										QuantityOfRegistersHight: 0,
										QuantityOfRegistersLow:   7,
									},
									ByteCount: 14,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
								},
								Response: &structs.RTUOverTCPMultipleWriteResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      16,
										ErrorCheckLow:   96,
										ErrorCheckHight: 5,
									},
									RegisterAddressHight:     0,
									RegisterAddressLow:       150,
									QuantityOfRegistersHight: 0,
									QuantityOfRegistersLow:   7,
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 613916674, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "17",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      3,
										ErrorCheckLow:   228,
										ErrorCheckHight: 0,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   150,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     15,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      3,
										ErrorCheckLow:   42,
										ErrorCheckHight: 116,
									},
									ByteCount: 30,
									Data:      []uint16{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 631440949, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "18",
							},
							Handshake: structs.Handshake{
								Request: &structs.RTUOverTCPRequest123456Response56{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      4,
										ErrorCheckLow:   48,
										ErrorCheckHight: 36,
									},
									StartingAddressHight: 0,
									StartingAddressLow:   4,
									ReadWriteDataHight:   0,
									ReadWriteDataLow:     18,
								},
								Response: &structs.RTUOverTCPReadResponse{
									HeaderError: structs.HeaderErrorCheck{
										SlaveAddress:    3,
										FunctionID:      4,
										ErrorCheckLow:   153,
										ErrorCheckHight: 69,
									},
									ByteCount: 36,
									Data:      []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 9, 58, 24, 641064525, time.Local),
						},
					},
					Slaves: []uint8{1, 2, 3},
				},
			},
		},
		{
			mode:          "tcp",
			directoryPath: `pcapng_files/tests_files/multiple_ports`,
			ports: map[uint16]conf.ServerSocket{
				1502: {
					HostAddress: "127.0.0.1",
					PortAddress: "1502",
				},
				1503: {
					HostAddress: "127.0.0.1",
					PortAddress: "1503",
				},
			},
			expectedHistory: map[uint16]structs.ServerHistory{
				1502: {
					Transactions: []structs.HistoryEvent{
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-1",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    10,
										UnitID:        1,
										FunctionType:  15,
									},
									AddressStart: []byte{0, 7},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 3},
										NumberBits:      3,
										Data:            []byte{1, 1, 0},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  15,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 7},
										NumberWrittenRegisters: []byte{0, 3},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 493372239, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-2",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  1,
									},
									AddressStart: []byte{0, 5},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 5},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    4,
										UnitID:        1,
										FunctionType:  1,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 1,
										Bits:       4,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 493478544, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-3",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  2,
									},
									AddressStart: []byte{0, 9},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 11},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    5,
										UnitID:        1,
										FunctionType:  2,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 2,
										Bits:       146,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493601380, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-4",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    21,
										UnitID:        1,
										FunctionType:  16,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 7},
										NumberBits:      14,
										Data:            []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  16,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 150},
										NumberWrittenRegisters: []byte{0, 7},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493684938, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-5",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  3,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 15},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    33,
										UnitID:        1,
										FunctionType:  3,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 30,
										Data:       []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493773451, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-6",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  4,
									},
									AddressStart: []byte{0, 4},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 18},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    39,
										UnitID:        1,
										FunctionType:  4,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 36,
										Data:       []byte{0, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 25, 248},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493891013, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-1",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    10,
										UnitID:        2,
										FunctionType:  15,
									},
									AddressStart: []byte{0, 7},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 3},
										NumberBits:      3,
										Data:            []byte{1, 1, 0},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  15,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 7},
										NumberWrittenRegisters: []byte{0, 3},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494213430, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-2",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  1,
									},
									AddressStart: []byte{0, 5},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 5},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    4,
										UnitID:        2,
										FunctionType:  1,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 1,
										Bits:       4,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494315312, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-3",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  2,
									},
									AddressStart: []byte{0, 9},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 11},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    5,
										UnitID:        2,
										FunctionType:  2,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 2,
										Bits:       0,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494387555, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-4",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    21,
										UnitID:        2,
										FunctionType:  16,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 7},
										NumberBits:      14,
										Data:            []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  16,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 150},
										NumberWrittenRegisters: []byte{0, 7},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494470536, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-5",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  3,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 15},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    33,
										UnitID:        2,
										FunctionType:  3,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 30,
										Data:       []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494546296, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-6",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  4,
									},
									AddressStart: []byte{0, 4},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 18},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    39,
										UnitID:        2,
										FunctionType:  4,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 36,
										Data:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494617661, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-1",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    10,
										UnitID:        3,
										FunctionType:  15,
									},
									AddressStart: []byte{0, 7},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 3},
										NumberBits:      3,
										Data:            []byte{1, 1, 0},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  15,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 7},
										NumberWrittenRegisters: []byte{0, 3},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494933890, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-2",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  1,
									},
									AddressStart: []byte{0, 5},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 5},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    4,
										UnitID:        3,
										FunctionType:  1,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 1,
										Bits:       4,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 495003356, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-3",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  2,
									},
									AddressStart: []byte{0, 9},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 11},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    5,
										UnitID:        3,
										FunctionType:  2,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 2,
										Bits:       0,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495068058, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-4",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    21,
										UnitID:        3,
										FunctionType:  16,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 7},
										NumberBits:      14,
										Data:            []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  16,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 150},
										NumberWrittenRegisters: []byte{0, 7},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495145227, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-5",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  3,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 15},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    33,
										UnitID:        3,
										FunctionType:  3,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 30,
										Data:       []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495220890, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-6",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  4,
									},
									AddressStart: []byte{0, 4},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 18},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    39,
										UnitID:        3,
										FunctionType:  4,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 36,
										Data:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495311180, time.Local),
						},
					},
					Slaves: []uint8{1, 2, 3},
				},
				1503: {
					Transactions: []structs.HistoryEvent{
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-1",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    10,
										UnitID:        1,
										FunctionType:  15,
									},
									AddressStart: []byte{0, 7},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 3},
										NumberBits:      3,
										Data:            []byte{1, 1, 0},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  15,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 7},
										NumberWrittenRegisters: []byte{0, 3},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 493482382, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-2",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  1,
									},
									AddressStart: []byte{0, 5},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 5},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    4,
										UnitID:        1,
										FunctionType:  1,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 1,
										Bits:       4,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 493582751, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-3",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  2,
									},
									AddressStart: []byte{0, 9},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 11},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    5,
										UnitID:        1,
										FunctionType:  2,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 2,
										Bits:       146,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493674285, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-4",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    21,
										UnitID:        1,
										FunctionType:  16,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 7},
										NumberBits:      14,
										Data:            []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  16,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 150},
										NumberWrittenRegisters: []byte{0, 7},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493762991, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-5",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  3,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 15},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    33,
										UnitID:        1,
										FunctionType:  3,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 30,
										Data:       []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493854239, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       1,
								TransactionID: "0-6",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        1,
										FunctionType:  4,
									},
									AddressStart: []byte{0, 4},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 18},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    39,
										UnitID:        1,
										FunctionType:  4,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 36,
										Data:       []byte{0, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 25, 248},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 493933699, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-1",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    10,
										UnitID:        2,
										FunctionType:  15,
									},
									AddressStart: []byte{0, 7},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 3},
										NumberBits:      3,
										Data:            []byte{1, 1, 0},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  15,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 7},
										NumberWrittenRegisters: []byte{0, 3},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494203449, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-2",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  1,
									},
									AddressStart: []byte{0, 5},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 5},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    4,
										UnitID:        2,
										FunctionType:  1,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 1,
										Bits:       4,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494285109, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-3",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  2,
									},
									AddressStart: []byte{0, 9},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 11},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    5,
										UnitID:        2,
										FunctionType:  2,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 2,
										Bits:       0,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494359143, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-4",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    21,
										UnitID:        2,
										FunctionType:  16,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 7},
										NumberBits:      14,
										Data:            []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  16,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 150},
										NumberWrittenRegisters: []byte{0, 7},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494426128, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-5",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  3,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 15},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    33,
										UnitID:        2,
										FunctionType:  3,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 30,
										Data:       []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494504852, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       2,
								TransactionID: "0-6",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        2,
										FunctionType:  4,
									},
									AddressStart: []byte{0, 4},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 18},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    39,
										UnitID:        2,
										FunctionType:  4,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 36,
										Data:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494599249, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-1",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    10,
										UnitID:        3,
										FunctionType:  15,
									},
									AddressStart: []byte{0, 7},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 3},
										NumberBits:      3,
										Data:            []byte{1, 1, 0},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 1},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  15,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 7},
										NumberWrittenRegisters: []byte{0, 3},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494801609, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-2",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  1,
									},
									AddressStart: []byte{0, 5},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 5},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 2},
										Protocol:      "modbus",
										BodyLength:    4,
										UnitID:        3,
										FunctionType:  1,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 1,
										Bits:       4,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 6, 10, 01, 21, 494901640, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-3",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  2,
									},
									AddressStart: []byte{0, 9},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 11},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 3},
										Protocol:      "modbus",
										BodyLength:    5,
										UnitID:        3,
										FunctionType:  2,
									},
									Data: &structs.TCPReadBitResponse{
										NumberBits: 2,
										Bits:       0,
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 494977332, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-4",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    21,
										UnitID:        3,
										FunctionType:  16,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPWriteMultipleRequest{
										NumberRegisters: []byte{0, 7},
										NumberBits:      14,
										Data:            []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 4},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  16,
									},
									Data: &structs.TCPWriteMultipleResponse{
										AddressStart:           []byte{0, 150},
										NumberWrittenRegisters: []byte{0, 7},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495055886, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-5",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  3,
									},
									AddressStart: []byte{0, 150},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 15},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 5},
										Protocol:      "modbus",
										BodyLength:    33,
										UnitID:        3,
										FunctionType:  3,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 30,
										Data:       []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495142910, time.Local),
						},
						{
							Header: structs.SlaveTransaction{
								SlaveID:       3,
								TransactionID: "0-6",
							},
							Handshake: structs.Handshake{
								Request: &structs.TCPRequest{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    6,
										UnitID:        3,
										FunctionType:  4,
									},
									AddressStart: []byte{0, 4},
									Data: &structs.TCPReadRequest{
										NumberReadingBits: []byte{0, 18},
									},
								},
								Response: &structs.TCPResponse{
									Header: structs.MBAPHeader{
										TransactionID: []byte{0, 6},
										Protocol:      "modbus",
										BodyLength:    39,
										UnitID:        3,
										FunctionType:  4,
									},
									Data: &structs.TCPReadByteResponse{
										NumberBits: 36,
										Data:       []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
									},
								},
							},
							TransactionTime: time.Date(2024, 12, 06, 10, 1, 21, 495242243, time.Local),
						},
					},
					Slaves: []uint8{1, 2, 3},
				},
			},
		},
	}
	var currentRecievedHistory map[uint16]structs.ServerHistory
	var err error
	for _, currentTestCase := range testTable {
		conf.WorkMode = currentTestCase.mode
		conf.DumpDirectoryPath = currentTestCase.directoryPath
		conf.Ports = currentTestCase.ports
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
