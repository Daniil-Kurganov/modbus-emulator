package tests_test

import (
	"io/ioutil"
	"log"
	"modbus-emulator/src"
	"modbus-emulator/src/utils"
	"testing"
	"time"

	mc "github.com/goburrow/modbus"
	"github.com/stretchr/testify/assert"
)

type (
	addresses struct {
		start    uint16
		quantity uint16
	}
	registersTCP struct {
		coils []byte
		DI    []byte
		HR    []byte
		IR    []byte
	}
	registersRTUOverTCP struct {
		coils []bool
		DI    []bool
		HR    []uint16
		IR    []uint16
	}
	transactionValues[T registersTCP | registersRTUOverTCP] struct {
		delay            time.Duration
		browsedRegisters map[string]addresses
		expectedStates   T
	}
	testCase[T registersTCP | registersRTUOverTCP] struct {
		workMode     string
		transactions []transactionValues[T]
	}
)

func TestServerTCPMode(t *testing.T) {
	var err error
	log.SetOutput(ioutil.Discard)
	directoryPath := `src/pcapng_files/tests_files`
	testCasesTCP := testCase[registersTCP]{
		workMode: "tcp",
		transactions: []transactionValues[registersTCP]{
			{
				delay: 501 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 0, quantity: 1},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersTCP{
					coils: []byte{0},
					DI:    []byte{0},
					HR:    []byte{0, 0},
					IR:    []byte{0},
				},
			},
			{
				delay: 903 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 0, quantity: 1},
					"DI":    {start: 16, quantity: 2},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersTCP{
					coils: []byte{0},
					DI:    []byte{0},
					HR:    []byte{0, 0},
					IR:    []byte{0},
				},
			},
			{
				delay: 1100 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 0, quantity: 1},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 8, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersTCP{
					coils: []byte{0},
					DI:    []byte{0},
					HR:    []byte{0, 39},
					IR:    []byte{0},
				},
			},
			{
				delay: 101 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 5, quantity: 5},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersTCP{
					coils: []byte{5},
					DI:    []byte{0},
					HR:    []byte{0, 0},
					IR:    []byte{0},
				},
			},
			{
				delay: 3 * time.Second,
				browsedRegisters: map[string]addresses{
					"coils": {start: 4, quantity: 4},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersTCP{
					coils: []byte{11},
					DI:    []byte{0},
					HR:    []byte{0, 0},
					IR:    []byte{0},
				},
			},
		},
	}
	utils.DumpDirectoryPath = directoryPath
	utils.WorkMode = testCasesTCP.workMode
	go src.ServerInit()
	handler := mc.NewTCPClientHandler("localhost:1502")
	if err = handler.Connect(); err != nil {
		log.Fatalf("Error on handler connecting: %s\n", err)
	}
	defer handler.Close()
	client := mc.NewClient(handler)
	for _, currentTCPTestCase := range testCasesTCP.transactions {
		var currentRecievedStates registersTCP
		if currentRecievedStates.coils, err = client.ReadCoils(currentTCPTestCase.browsedRegisters["coils"].start,
			currentTCPTestCase.browsedRegisters["coils"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		if currentRecievedStates.DI, err = client.ReadDiscreteInputs(currentTCPTestCase.browsedRegisters["DI"].start,
			currentTCPTestCase.browsedRegisters["DI"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		if currentRecievedStates.HR, err = client.ReadHoldingRegisters(currentTCPTestCase.browsedRegisters["HR"].start,
			currentTCPTestCase.browsedRegisters["HR"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		if currentRecievedStates.IR, err = client.ReadDiscreteInputs(currentTCPTestCase.browsedRegisters["IR"].start,
			currentTCPTestCase.browsedRegisters["IR"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		assert.Equalf(t, currentTCPTestCase.expectedStates, currentRecievedStates,
			"Error: recieved and expected states isn't equal:\n expected: %v;\n recieved: %v",
			currentTCPTestCase.expectedStates, currentRecievedStates)
		time.Sleep(currentTCPTestCase.delay)
	}
}
