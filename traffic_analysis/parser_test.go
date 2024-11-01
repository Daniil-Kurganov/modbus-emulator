package trafficanalysis_test

import (
	"fmt"
	ta "modbus-emulator/traffic_analysis"
	"modbus-emulator/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePackets(t *testing.T) {
	testTable := []struct {
		filename         string
		responseExpected []ta.TCPPacket
		errorExpected    error
	}{
		{
			filename:         "",
			responseExpected: nil,
			errorExpected:    fmt.Errorf("error on opening file: %s/%s/.pcapng: No such file or directory", utils.ModulePath, utils.Foldername),
		},
		{
			filename: "coils_read",
			responseExpected: []ta.TCPPacket{
				{
					PacketNumber: 1,
					Protocol:     "modbus",
					BodyLength:   4,
					UnitID:       0,
					ObjectType:   1,
					DataLength:   1,
					Data:         []byte{1},
				},
				{
					PacketNumber: 2,
					Protocol:     "modbus",
					BodyLength:   4,
					UnitID:       0,
					ObjectType:   1,
					DataLength:   1,
					Data:         []byte{1},
				},
				{
					PacketNumber: 3,
					Protocol:     "modbus",
					BodyLength:   4,
					UnitID:       0,
					ObjectType:   1,
					DataLength:   1,
					Data:         []byte{1},
				},
				{
					PacketNumber: 4,
					Protocol:     "modbus",
					BodyLength:   4,
					UnitID:       0,
					ObjectType:   1,
					DataLength:   1,
					Data:         []byte{1},
				},
				{
					PacketNumber: 5,
					Protocol:     "modbus",
					BodyLength:   4,
					UnitID:       0,
					ObjectType:   1,
					DataLength:   1,
					Data:         []byte{1},
				},
			},
			errorExpected: nil,
		},
	}
	for _, currentTestCase := range testTable {
		var currentResponseRecieved []ta.TCPPacket
		var err error
		if currentResponseRecieved, err = ta.ParsePackets(currentTestCase.filename); err != nil {
			assert.EqualErrorf(t, err, currentTestCase.errorExpected.Error(),
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", currentTestCase.errorExpected, err,
			)
		}
		assert.Equal(t, currentTestCase.responseExpected, currentResponseRecieved,
			fmt.Sprintf("Error: recieved and expected responses isn't equal:\n expected: %v;\n recieved: %v",
				currentTestCase.responseExpected, currentResponseRecieved),
		)
	}
}
