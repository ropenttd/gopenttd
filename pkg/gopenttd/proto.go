// Package gopenttd provides primitives for querying OpenTTD game servers.
package gopenttd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skybon/goutil"
)

func getByteString(byteArray []byte) string {
	return fmt.Sprintf("%x", byteArray)
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
// Connections time out after 10 seconds
func ScanServer(host string, port int) (serverstate OpenttdServerState, err error) {
	var timeout = time.Second * 10

	inBuf := make([]byte, 1024)
	var readLen int

	// Determine the correct UDP address
	server := fmt.Sprintf("%s:%d", host, port) // "255.255.255.255", 10000
	serverAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return OpenttdServerState{Status: false, Error: err}, err
	}

	// Open the conneciton
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return OpenttdServerState{Status: false, Error: err}, err
	}
	log.Debug("Connection open to ", server)
	defer conn.Close()
	// Set a timeout on the connection
	conn.SetDeadline(time.Now().Add(timeout))

	// send the Info request
	discoverMsg := []byte("\x03\x00\x00")
	sendLen, err := conn.Write(discoverMsg)
	if err != nil {
		return OpenttdServerState{Status: false, Error: err}, err
	}
	log.Debug("Sent ", sendLen, " bytes")

	for {
		// Start a loop to read from UDP, and exit out if an error or timeout is experienced
		readLen, _, err = conn.ReadFromUDP(inBuf)

		if err != nil {
			// Some kind of connection error
			return OpenttdServerState{Status: false, Error: err}, err
		}

		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			// Timeout error
			return OpenttdServerState{Status: false, Error: err}, err
		}

		if readLen != 0 {
			// Success!
			log.Debug("Read ", readLen, " bytes")
			break
		}
	}
	serverstate = OpenttdServerState{Host: serverAddr.String()}
	serverstate.Populate(inBuf)
	log.Fatalf("%x", inBuf)
	return
}

// Populate populates an OpenttdServerState struct with data parsed from the given Byte Array.
func (server *OpenttdServerState) Populate(p []byte) {
	// Props to https://github.com/vorot93/grokstat/blob/master/protocol_openttds.go for most of this code
	var infoData = bytes.NewBuffer(p[3:])

	var protocolVer = int(infoData.Next(1)[0])

	var activeNewGRFsNum int
	var activeNewGRFsInfo []OpenttdNewgrf
	if protocolVer >= 4 {
		activeNewGRFsNum = int(infoData.Next(1)[0])
		for n := 0; n < activeNewGRFsNum; n++ {
			NewGRFID := getByteString(infoData.Next(4))
			NewGRFMD5 := getByteString(infoData.Next(16))
			activeNewGRFsInfo = append(activeNewGRFsInfo, OpenttdNewgrf{Identifier: NewGRFID, Hash: NewGRFMD5})
		}
	}

	var timeCurrent time.Time
	var timeStart time.Time
	if protocolVer >= 3 {
		var tc uint32
		var ts uint32
		_ = binary.Read(bytes.NewReader(infoData.Next(4)), binary.LittleEndian, &tc)
		_ = binary.Read(bytes.NewReader(infoData.Next(4)), binary.LittleEndian, &ts)
		timeCurrent = OttdDateFormat(tc)
		timeStart = OttdDateFormat(ts)
	}

	var maxCompanies *int
	var currentCompanies *int
	var maxSpectators *int
	if protocolVer >= 2 {
		maxCompanies = goutil.IntP(int(infoData.Next(1)[0]))
		currentCompanies = goutil.IntP(int(infoData.Next(1)[0]))
		maxSpectators = goutil.IntP(int(infoData.Next(1)[0]))
	}
	serverNameBytes, _ := infoData.ReadBytes(byte(0))
	serverName := string(bytes.Trim(serverNameBytes, "\x00"))

	serverVersionBytes, _ := infoData.ReadBytes(byte(0))
	serverVersion := string(bytes.Trim(serverVersionBytes, "\x00"))

	languageId := int(infoData.Next(1)[0])
	needPass := bool(infoData.Next(1)[0] != 0)
	maxClients := int(infoData.Next(1)[0])
	currentClients := int(infoData.Next(1)[0])
	currentSpectators := int(infoData.Next(1)[0])

	if protocolVer < 3 {
		_ = infoData.Next(2)
		_ = infoData.Next(2)
	}

	mapNameBytes, _ := infoData.ReadBytes(byte(0))
	mapName := string(bytes.Trim(mapNameBytes, "\x00"))

	var mapWidth uint16
	_ = binary.Read(bytes.NewReader(infoData.Next(2)), binary.LittleEndian, &mapWidth)

	var mapHeight uint16
	_ = binary.Read(bytes.NewReader(infoData.Next(2)), binary.LittleEndian, &mapHeight)

	mapSet := int(infoData.Next(1)[0])
	dedicatedServer := int(infoData.Next(1)[0])

	server.Status = true
	server.Dedicated = !(dedicatedServer == 0)
	server.Name = fmt.Sprint(serverName)
	server.Version = fmt.Sprint(serverVersion)
	server.Language = languageId
	server.NeedPass = needPass
	server.Environment = mapSet
	server.Map = mapName
	server.MapWidth = mapWidth
	server.MapHeight = mapHeight
	server.DateStart = timeStart
	server.DateCurrent = timeCurrent
	server.NumClients = currentClients
	server.MaxClients = maxClients
	server.NumSpectators = currentSpectators
	server.MaxSpectators = *maxSpectators
	server.NumCompanies = *currentCompanies
	server.MaxCompanies = *maxCompanies
	server.NewgrfCount = activeNewGRFsNum
	server.NewgrfActive = activeNewGRFsInfo

}
