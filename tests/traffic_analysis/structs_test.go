package trafficanalysis_test

import (
	"testing"
	ta "modbus-emulator/src/traffic_analysis"
)

func TestMBAPHeaderUnmarshal(t *testing.T) {
	testTable := []struct {
		payload []byte
		expectedHeader ta.MBAPHeader
	}{
		{
			payload: []byte{0, 23, 0, 0, 0, 6, 0, 3, 0, 3, 0, 6},
			expectedHeader: ta.MBAPHeader{

			},
		},
	}
}