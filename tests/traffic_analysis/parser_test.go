package trafficanalysis_test

import (
	"fmt"
	"modbus-emulator/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	ta "modbus-emulator/src/traffic_analysis"
)

func TestParsePackets(t *testing.T) {
	testTable := []struct {
		typeObject string
		filename string
		expectedHistory []map[string]ta.Handshake
	}{
		{
			typeObject: "coils",
			filename: "read_01",
			expectedHistory: map[string]ta.Handshake{
				""
			},
		}
	}
	// for _, currentTestCase := range testTable {
	// 	var currentResponseRecieved []ta.TCPPacket
	// 	var err error
	// 	if currentResponseRecieved, err = ta.ParsePackets(currentTestCase.typeObject, currentTestCase.filename, currentTestCase.filter); err != nil {
	// 		assert.EqualErrorf(t, err, currentTestCase.errorExpected.Error(),
	// 			"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", currentTestCase.errorExpected, err,
	// 		)
	// 	}
	// 	assert.Equal(t, currentTestCase.responseExpected, currentResponseRecieved,
	// 		fmt.Sprintf("Error: recieved and expected responses isn't equal:\n expected: %v;\n recieved: %v",
	// 			currentTestCase.responseExpected, currentResponseRecieved),
	// 	)
	// }
}
