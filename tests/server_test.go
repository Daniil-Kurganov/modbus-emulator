package tests_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"modbus-emulator/conf"
	"modbus-emulator/src"
	"sync"
	"testing"
	"time"

	mc "github.com/goburrow/modbus"
	"github.com/simonvetter/modbus"
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

var (
	directoryPath = `pcapng_files/tests_files/simple_port`
	port          = 1502
)

func TestServerTCPMode(t *testing.T) {
	var err error
	log.SetOutput(ioutil.Discard)
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
					IR:    []byte{0, 0},
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
					IR:    []byte{0, 0},
				},
			},
			{
				delay: 1300 * time.Millisecond,
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
					IR:    []byte{0, 0},
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
					IR:    []byte{0, 0},
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
					IR:    []byte{0, 0},
				},
			},
		},
	}
	conf.DumpDirectoryPath = directoryPath
	conf.WorkMode = testCasesTCP.workMode
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go src.ServerInit(&waitGroup, uint16(port))
	time.Sleep(500 * time.Millisecond)
	handler := mc.NewTCPClientHandler(fmt.Sprintf("%s:%d", conf.ServerTCPHost, uint16(port)))
	handler.SlaveId = 0
	if err = handler.Connect(); err != nil {
		assert.EqualErrorf(t, err, "nil",
			"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
		)
		return
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
			return
		}
		if currentRecievedStates.DI, err = client.ReadDiscreteInputs(currentTCPTestCase.browsedRegisters["DI"].start,
			currentTCPTestCase.browsedRegisters["DI"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
			return
		}
		if currentRecievedStates.HR, err = client.ReadHoldingRegisters(currentTCPTestCase.browsedRegisters["HR"].start,
			currentTCPTestCase.browsedRegisters["HR"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
			return
		}
		if currentRecievedStates.IR, err = client.ReadInputRegisters(currentTCPTestCase.browsedRegisters["IR"].start,
			currentTCPTestCase.browsedRegisters["IR"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
			return
		}
		assert.Equalf(t, currentTCPTestCase.expectedStates, currentRecievedStates,
			"Error: recieved and expected states isn't equal:\n expected: %v;\n recieved: %v",
			currentTCPTestCase.expectedStates, currentRecievedStates)
		time.Sleep(currentTCPTestCase.delay)
	}
	waitGroup.Wait()
}

func TestServerRTUOverTCPMode(t *testing.T) {
	var err error
	log.SetOutput(ioutil.Discard)
	testCasesRTUOverTCP := testCase[registersRTUOverTCP]{
		workMode: "rtu_over_tcp",
		transactions: []transactionValues[registersRTUOverTCP]{
			{
				delay: 507 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 0, quantity: 1},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersRTUOverTCP{
					coils: []bool{false},
					DI:    []bool{false},
					HR:    []uint16{0},
					IR:    []uint16{0},
				},
			},
			{
				delay: 909 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 0, quantity: 1},
					"DI":    {start: 0, quantity: 2},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersRTUOverTCP{
					coils: []bool{false},
					DI:    []bool{false, false},
					HR:    []uint16{0},
					IR:    []uint16{0},
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
				expectedStates: registersRTUOverTCP{
					coils: []bool{false},
					DI:    []bool{false},
					HR:    []uint16{39},
					IR:    []uint16{0},
				},
			},
			{
				delay: 107 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 5, quantity: 5},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersRTUOverTCP{
					coils: []bool{true, false, true, false, false},
					DI:    []bool{false},
					HR:    []uint16{0},
					IR:    []uint16{0},
				},
			},
			{
				delay: 1100 * time.Millisecond,
				browsedRegisters: map[string]addresses{
					"coils": {start: 4, quantity: 4},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 0, quantity: 1},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersRTUOverTCP{
					coils: []bool{true, true, false, true},
					DI:    []bool{false},
					HR:    []uint16{0},
					IR:    []uint16{0},
				},
			},
			{
				delay: 3 * time.Second,
				browsedRegisters: map[string]addresses{
					"coils": {start: 0, quantity: 1},
					"DI":    {start: 0, quantity: 1},
					"HR":    {start: 3, quantity: 4},
					"IR":    {start: 0, quantity: 1},
				},
				expectedStates: registersRTUOverTCP{
					coils: []bool{false},
					DI:    []bool{false},
					HR:    []uint16{0, 16, 0, 0},
					IR:    []uint16{0},
				},
			},
		},
	}
	conf.DumpDirectoryPath = directoryPath
	conf.WorkMode = testCasesRTUOverTCP.workMode
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go src.ServerInit(&waitGroup, uint16(port))
	time.Sleep(500 * time.Millisecond)
	var client *modbus.ModbusClient
	if client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("rtuovertcp://%s:%d", conf.ServerTCPHost, uint16(port)),
		Speed:   19200,
		Timeout: 1 * time.Second,
	}); err != nil {
		assert.EqualErrorf(t, err, "nil",
			"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
		)
		return
	}
	client.SetUnitId(1)
	if err = client.Open(); err != nil {
		assert.EqualErrorf(t, err, "nil",
			"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
		)
		return
	}
	defer client.Close()
	for _, currentRTUOverTCPTestCase := range testCasesRTUOverTCP.transactions {
		var currentRecievedStates registersRTUOverTCP
		if currentRecievedStates.coils, err = client.ReadCoils(currentRTUOverTCPTestCase.browsedRegisters["coils"].start,
			currentRTUOverTCPTestCase.browsedRegisters["coils"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		if currentRecievedStates.DI, err = client.ReadDiscreteInputs(currentRTUOverTCPTestCase.browsedRegisters["DI"].start,
			currentRTUOverTCPTestCase.browsedRegisters["DI"].quantity); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		if currentRecievedStates.HR, err = client.ReadRegisters(currentRTUOverTCPTestCase.browsedRegisters["HR"].start,
			currentRTUOverTCPTestCase.browsedRegisters["HR"].quantity, modbus.HOLDING_REGISTER); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		if currentRecievedStates.IR, err = client.ReadRegisters(currentRTUOverTCPTestCase.browsedRegisters["IR"].start,
			currentRTUOverTCPTestCase.browsedRegisters["IR"].quantity, modbus.INPUT_REGISTER); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
		}
		assert.Equalf(t, currentRTUOverTCPTestCase.expectedStates, currentRecievedStates,
			"Error: recieved and expected states isn't equal:\n expected: %v;\n recieved: %v",
			currentRTUOverTCPTestCase.expectedStates, currentRecievedStates)
		time.Sleep(currentRTUOverTCPTestCase.delay)
	}
	waitGroup.Wait()
}
