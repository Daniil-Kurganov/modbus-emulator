package tests_test

import (
	"modbus-emulator/src/traffic_analysis/structs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmulationData(t *testing.T) {
	testCases := []struct {
		handshake             structs.Handshake
		expectedEmulationData structs.EmulationData
	}{
		{
			structs.Handshake{
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
			structs.EmulationData{
				FunctionID:      3,
				IsReadOperation: true,
				Address:         150,
				Quantity:        15,
				Payload:         []uint16{1, 18, 48, 53, 64, 57, 59, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
		{
			structs.Handshake{
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
						BodyLength:    10,
						UnitID:        3,
						FunctionType:  15,
					},
					Data: &structs.TCPWriteMultipleResponse{
						AddressStart:           []byte{0, 7},
						NumberWrittenRegisters: []byte{0, 3},
					},
				},
			},
			structs.EmulationData{
				FunctionID:      15,
				IsReadOperation: false,
				Address:         7,
				Quantity:        3,
				Payload:         []uint16{1, 1, 0},
			},
		},
	}
	for _, currentTestCase := range testCases {
		var currentRecievedEmulationData structs.EmulationData
		var err error
		if currentRecievedEmulationData, err = currentTestCase.handshake.Marshal(); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		assert.Equalf(t, currentTestCase.expectedEmulationData, currentRecievedEmulationData,
			"Error: recieved and expected emulations data isn't equal:\n expected: %+v;\n recieved: %+v",
			currentTestCase.expectedEmulationData, currentRecievedEmulationData,
		)
	}
}
