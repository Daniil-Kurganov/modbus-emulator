package tests__test

import (
	"io/ioutil"
	"log"
	"modbus-emulator/src"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	addresses struct {
		start int
		end   int
	}
	testCase struct {
		delay            time.Duration
		browsedRegisters map[string]addresses
		expectedStates   map[string][]byte
	}
)

func TestServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	testCases := []testCase{
		{
			delay: 501 * time.Millisecond,
			browsedRegisters: map[string]addresses{
				"coils": {start: 0, end: 1},
				"DI":    {start: 0, end: 1},
				"HR":    {start: 0, end: 1},
				"IR":    {start: 0, end: 1},
			},
			expectedStates: map[string][]byte{
				"coils": {0},
				"DI":    {0},
				"HR":    {0},
				"IR":    {0},
			},
		},
		{
			delay: 903 * time.Millisecond,
			browsedRegisters: map[string]addresses{
				"coils": {start: 0, end: 1},
				"DI":    {start: 16, end: 18},
				"HR":    {start: 0, end: 1},
				"IR":    {start: 0, end: 1},
			},
			expectedStates: map[string][]byte{
				"coils": {0},
				"DI":    {0, 0},
				"HR":    {0},
				"IR":    {0},
			},
		},
		{
			delay: 1100 * time.Millisecond,
			browsedRegisters: map[string]addresses{
				"coils": {start: 0, end: 1},
				"DI":    {start: 0, end: 1},
				"HR":    {start: 8, end: 9},
				"IR":    {start: 0, end: 1},
			},
			expectedStates: map[string][]byte{
				"coils": {0},
				"DI":    {0},
				"HR":    {0},
				"IR":    {0},
			},
		},
		{
			delay: 101 * time.Millisecond,
			browsedRegisters: map[string]addresses{
				"coils": {start: 5, end: 10},
				"DI":    {start: 0, end: 1},
				"HR":    {start: 0, end: 1},
				"IR":    {start: 0, end: 1},
			},
			expectedStates: map[string][]byte{
				"coils": {1, 0, 1, 0, 0},
				"DI":    {0},
				"HR":    {0},
				"IR":    {0},
			},
		},
		{
			delay: 3 * time.Second,
			browsedRegisters: map[string]addresses{
				"coils": {start: 4, end: 8},
				"DI":    {start: 0, end: 1},
				"HR":    {start: 0, end: 1},
				"IR":    {start: 0, end: 1},
			},
			expectedStates: map[string][]byte{
				"coils": {1, 1, 0, 1},
				"DI":    {0},
				"HR":    {0},
				"IR":    {0},
			},
		},
	}
	src.ServerInit()
	for _, currentTestCase := range testCases {
		currentRecievedStates := map[string][]byte{
			"coils": src.Server.Coils[currentTestCase.browsedRegisters["coils"].start:currentTestCase.browsedRegisters["coils"].end],
			"DI":    src.Server.Coils[currentTestCase.browsedRegisters["DI"].start:currentTestCase.browsedRegisters["DI"].end],
			"HR":    src.Server.Coils[currentTestCase.browsedRegisters["HR"].start:currentTestCase.browsedRegisters["HR"].end],
			"IR":    src.Server.Coils[currentTestCase.browsedRegisters["IR"].start:currentTestCase.browsedRegisters["IR"].end],
		}
		assert.Equalf(t, currentTestCase.expectedStates, currentRecievedStates,
			"Error: recieved and expected states isn't equal:\n expected: %v;\n recieved: %v",
			currentTestCase.expectedStates, currentRecievedStates)
	}
}
