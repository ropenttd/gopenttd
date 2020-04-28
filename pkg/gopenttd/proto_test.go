package gopenttd

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDecodePacketNoGRF(t *testing.T) {
	testPacket, _ := hex.DecodeString("4f0001040028380b00dacf0a000f0b2d7265646469742e636f6d2f722f4f70656e545444202331202d2056616e696c6c6100312e31302e3100000032060052616e646f6d204d617000000400080001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	expectObject := OpenttdServerState{
		Status:        true,
		Error:         nil,
		Dedicated:     true,
		Name:          "reddit.com/r/OpenTTD #1 - Vanilla",
		Version:       "1.10.1",
		NeedPass:      false,
		Language:      0,
		Environment:   0,
		Map:           "Random Map",
		MapHeight:     0,
		MapWidth:      0,
		DateStart:     time.Date(1940, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		DateCurrent:   time.Date(2013, time.Month(2), 8, 0, 0, 0, 0, time.UTC),
		NumClients:    6,
		MaxClients:    50,
		NumSpectators: 0,
		MaxSpectators: 45,
		NumCompanies:  11,
		MaxCompanies:  15,
		NewgrfCount:   0,
		NewgrfActive:  nil,
	}

	output := OpenttdServerState{}
	output.Populate(testPacket)

	assert.EqualValues(t, expectObject, output)
}
