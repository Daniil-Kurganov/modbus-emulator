package tests_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"modbus-emulator/conf"
	"modbus-emulator/src"
	ta "modbus-emulator/src/traffic_analysis"
	structs "modbus-emulator/src/traffic_analysis/structs"
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
		transactions map[uint8][]transactionValues[T]
	}
)

func TestServerTCPMode(t *testing.T) {
	var err error
	log.SetOutput(ioutil.Discard)
	directoryPath := `pcapng_files/tests_files/simple_port`
	testCasesTCP := testCase[registersTCP]{
		workMode: "tcp",
		transactions: map[uint8][]transactionValues[registersTCP]{
			0: {
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
		},
	}
	conf.DumpFilePath = fmt.Sprintf("/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator/%s/%s",
		directoryPath, testCasesTCP.workMode,
	)
	conf.Sockets = map[string]conf.ServerSocketData{
		"127.0.0.1:1502": {
			HostAddress: "localhost",
			PortAddress: "1502",
			Protocol:    testCasesTCP.workMode,
		},
	}
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go src.ServerInit(&waitGroup, "127.0.0.1:1502", history["127.0.0.1:1502"])
	time.Sleep(500 * time.Millisecond)
	handler := mc.NewTCPClientHandler("127.0.0.1:1502")
	handler.SlaveId = 0
	if err = handler.Connect(); err != nil {
		assert.EqualErrorf(t, err, "nil",
			"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
		)
		return
	}
	defer handler.Close()
	client := mc.NewClient(handler)
	for _, currentTCPTestCase := range testCasesTCP.transactions[0] {
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
	directoryPath := `pcapng_files/tests_files/simple_port`
	socket := "127.0.0.1:1502"
	testCasesRTUOverTCP := testCase[registersRTUOverTCP]{
		workMode: "rtu_over_tcp",
		transactions: map[uint8][]transactionValues[registersRTUOverTCP]{
			1: {
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
		},
	}
	conf.DumpFilePath = fmt.Sprintf("/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator/%s/%s",
		directoryPath, testCasesRTUOverTCP.workMode,
	)
	conf.Sockets = map[string]conf.ServerSocketData{
		"127.0.0.1:1502": {
			HostAddress: "localhost",
			PortAddress: "1502",
			Protocol:    testCasesRTUOverTCP.workMode,
		},
	}
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go src.ServerInit(&waitGroup, socket, history[socket])
	time.Sleep(500 * time.Millisecond)
	var client *modbus.ModbusClient
	if client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("rtuovertcp://%s", socket),
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
	for _, currentRTUOverTCPTestCase := range testCasesRTUOverTCP.transactions[1] {
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

func TestServerRTUOverTCPMupliplePorts(t *testing.T) {
	var err error
	log.SetOutput(ioutil.Discard)
	testCases := map[string]testCase[registersRTUOverTCP]{
		"127.0.0.1:1502": {
			transactions: map[uint8][]transactionValues[registersRTUOverTCP]{
				1: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 10, quantity: 7},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 4, quantity: 18},
						},
						expectedStates: registersRTUOverTCP{
							coils: []bool{false, true, true},
							DI:    []bool{true, false, false, true, false, false, true},
							HR:    []uint16{1, 18, 48, 53, 64, 57, 59},
							IR:    []uint16{120, 0, 0, 0, 0, 0, 385, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 6648},
						},
					},
				},
				2: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersRTUOverTCP{
							coils: []bool{false, true, true},
							DI:    []bool{false},
							HR:    []uint16{1, 18, 48, 53, 64, 57, 59},
							IR:    []uint16{0},
						},
					},
				},
				3: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersRTUOverTCP{
							coils: []bool{false, true, true},
							DI:    []bool{false},
							HR:    []uint16{1, 18, 48, 53, 64, 57, 59},
							IR:    []uint16{0},
						},
					},
				},
			},
		},
		"127.0.0.1:1503": {
			transactions: map[uint8][]transactionValues[registersRTUOverTCP]{
				1: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 10, quantity: 7},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 4, quantity: 18},
						},
						expectedStates: registersRTUOverTCP{
							coils: []bool{false, true, true},
							DI:    []bool{true, false, false, true, false, false, true},
							HR:    []uint16{1, 18, 48, 53, 64, 57, 59},
							IR:    []uint16{120, 0, 0, 0, 0, 0, 385, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 6648},
						},
					},
				},
				2: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersRTUOverTCP{
							coils: []bool{false, true, true},
							DI:    []bool{false},
							HR:    []uint16{1, 18, 48, 53, 64, 57, 59},
							IR:    []uint16{0},
						},
					},
				},
				3: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersRTUOverTCP{
							coils: []bool{false, true, true},
							DI:    []bool{false},
							HR:    []uint16{1, 18, 48, 53, 64, 57, 59},
							IR:    []uint16{0},
						},
					},
				},
			},
		},
	}
	conf.DumpFilePath = fmt.Sprintf("/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator/pcapng_files/tests_files/multiple_ports/%s",
		conf.Protocols.RTUOverTCP,
	)
	conf.Sockets = map[string]conf.ServerSocketData{
		"127.0.0.1:1502": {
			HostAddress: "localhost",
			PortAddress: "1502",
			Protocol:    conf.Protocols.RTUOverTCP,
		},
		"127.0.0.1:1503": {
			HostAddress: "localhost",
			PortAddress: "1503",
			Protocol:    conf.Protocols.RTUOverTCP,
		},
	}
	conf.FinishDelayTime = 3 * time.Second
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(testCases))
	for currentSocket, currentTestCase := range testCases {
		go src.ServerInit(&waitGroup, currentSocket, history[currentSocket])
		time.Sleep(500 * time.Millisecond)
		var client *modbus.ModbusClient
		if client, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL:     fmt.Sprintf("rtuovertcp://%s", currentSocket),
			Speed:   19200,
			Timeout: 1 * time.Second,
		}); err != nil {
			assert.EqualErrorf(t, err, "nil",
				"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
			)
			t.FailNow()
		}
		for currentSlaveId, currentTranscationValues := range currentTestCase.transactions {
			client.SetUnitId(currentSlaveId)
			if err = client.Open(); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentSlaveId == 1 {
				time.Sleep(1600 * time.Millisecond)
			}
			var currentRecievedStates registersRTUOverTCP
			if currentRecievedStates.coils, err = client.ReadCoils(currentTranscationValues[0].browsedRegisters["coils"].start,
				currentTranscationValues[0].browsedRegisters["coils"].quantity); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentRecievedStates.DI, err = client.ReadDiscreteInputs(currentTranscationValues[0].browsedRegisters["DI"].start,
				currentTranscationValues[0].browsedRegisters["DI"].quantity); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentRecievedStates.HR, err = client.ReadRegisters(currentTranscationValues[0].browsedRegisters["HR"].start,
				currentTranscationValues[0].browsedRegisters["HR"].quantity, modbus.HOLDING_REGISTER); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentRecievedStates.IR, err = client.ReadRegisters(currentTranscationValues[0].browsedRegisters["IR"].start,
				currentTranscationValues[0].browsedRegisters["IR"].quantity, modbus.INPUT_REGISTER); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if !assert.Equalf(t, currentTranscationValues[0].expectedStates, currentRecievedStates,
				"Error: recieved and expected states isn't equal:\n expected: %v;\n recieved: %v",
				currentTranscationValues[0].expectedStates, currentRecievedStates,
			) {
				t.FailNow()
			}
		}
		client.Close()
	}
	waitGroup.Wait()
}

func TestServerTCPMupliplePorts(t *testing.T) {
	var err error
	log.SetOutput(ioutil.Discard)
	testCases := map[string]testCase[registersTCP]{
		"127.0.0.1:1502": {
			transactions: map[uint8][]transactionValues[registersTCP]{
				1: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 10, quantity: 7},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 4, quantity: 18},
						},
						expectedStates: registersTCP{
							coils: []byte{2},
							DI:    []byte{73},
							HR:    []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
							IR:    []byte{0, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 25, 248},
						},
					},
				},
				2: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersTCP{
							coils: []byte{2},
							DI:    []byte{0},
							HR:    []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
							IR:    []byte{0, 0},
						},
					},
				},
				3: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersTCP{
							coils: []byte{2},
							DI:    []byte{0},
							HR:    []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
							IR:    []byte{0, 0},
						},
					},
				},
			},
		},
		"127.0.0.1:1503": {
			transactions: map[uint8][]transactionValues[registersTCP]{
				1: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 10, quantity: 7},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 4, quantity: 18},
						},
						expectedStates: registersTCP{
							coils: []byte{2},
							DI:    []byte{73},
							HR:    []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
							IR:    []byte{0, 120, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 129, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 0, 25, 248},
						},
					},
				},
				2: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersTCP{
							coils: []byte{2},
							DI:    []byte{0},
							HR:    []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
							IR:    []byte{0, 0},
						},
					},
				},
				3: {
					{
						browsedRegisters: map[string]addresses{
							"coils": {start: 6, quantity: 3},
							"DI":    {start: 1, quantity: 1},
							"HR":    {start: 150, quantity: 7},
							"IR":    {start: 0, quantity: 1},
						},
						expectedStates: registersTCP{
							coils: []byte{2},
							DI:    []byte{0},
							HR:    []byte{0, 1, 0, 18, 0, 48, 0, 53, 0, 64, 0, 57, 0, 59},
							IR:    []byte{0, 0},
						},
					},
				},
			},
		},
	}
	conf.DumpFilePath = fmt.Sprintf("/media/ugpa/1TB/Lavoro/Repositories/modbus-emulator/pcapng_files/tests_files/multiple_ports/%s",
		conf.Protocols.TCP,
	)
	conf.Sockets = map[string]conf.ServerSocketData{
		"127.0.0.1:1502": {
			HostAddress: "localhost",
			PortAddress: "1502",
			Protocol:    conf.Protocols.TCP,
		},
		"127.0.0.1:1503": {
			HostAddress: "localhost",
			PortAddress: "1503",
			Protocol:    conf.Protocols.TCP,
		},
	}
	conf.FinishDelayTime = 5 * time.Second
	var history map[string]structs.ServerHistory
	if history, err = ta.ParseDump(); err != nil {
		log.Fatalf("Error on parsing dump: %s", err)
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(testCases))
	for currentSocket, currentTestCase := range testCases {
		go src.ServerInit(&waitGroup, currentSocket, history[currentSocket])
		time.Sleep(1000 * time.Millisecond)
		handler := mc.NewTCPClientHandler(currentSocket)
		for currentSlaveId, currentTranscationValues := range currentTestCase.transactions {
			handler.SlaveId = currentSlaveId
			if err = handler.Connect(); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			client := mc.NewClient(handler)
			// if currentSlaveId == 1 {
			time.Sleep(1600 * time.Millisecond)
			// }
			var currentRecievedStates registersTCP
			if currentRecievedStates.coils, err = client.ReadCoils(currentTranscationValues[0].browsedRegisters["coils"].start,
				currentTranscationValues[0].browsedRegisters["coils"].quantity); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentRecievedStates.DI, err = client.ReadDiscreteInputs(currentTranscationValues[0].browsedRegisters["DI"].start,
				currentTranscationValues[0].browsedRegisters["DI"].quantity); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentRecievedStates.HR, err = client.ReadHoldingRegisters(currentTranscationValues[0].browsedRegisters["HR"].start,
				currentTranscationValues[0].browsedRegisters["HR"].quantity); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if currentRecievedStates.IR, err = client.ReadInputRegisters(currentTranscationValues[0].browsedRegisters["IR"].start,
				currentTranscationValues[0].browsedRegisters["IR"].quantity); err != nil {
				assert.EqualErrorf(t, err, "nil",
					"Error: recieved and expected errors isn't equal:\n expected: %s;\n recieved: %s", "nil", err,
				)
				t.FailNow()
			}
			if !assert.Equalf(t, currentTranscationValues[0].expectedStates, currentRecievedStates,
				"Error: recieved and expected states isn't equal:\n expected: %+v;\n recieved: %+v",
				currentTranscationValues[0].expectedStates, currentRecievedStates) {
				t.FailNow()
			}
			handler.Close()
		}
	}
	waitGroup.Wait()
}
