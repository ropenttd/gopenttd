// Package gopenttd provides primitives for querying OpenTTD game servers.
package gopenttd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/internal/openttd_packets_udp"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skybon/goutil"
)

func getByteString(byteArray []byte) string {
	return fmt.Sprintf("%x", byteArray)
}

func (c *OpenttdClientConnection) Open() (err error) {
	// Determine the correct UDP address
	server := fmt.Sprintf("%s:%d", c.Hostname, c.Port) // "255.255.255.255", 10000
	serverAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return err
	}

	// Open the connection
	c.Connection, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return err
	}
	log.Debug("Connection open to ", server)
	return nil
}

func (c *OpenttdClientConnection) Close() (err error) {
	return c.Connection.Close()
}

func (c *OpenttdClientConnection) Query(packetType uint8, expect uint8) (data *bytes.Buffer, err error) {
	// this buffer _should_ be long enough? probably a better way of doing this than allocating a 2m bytearray
	inBuf := make([]byte, 2048)
	var readLen int

	// Set a timeout
	var timeout = time.Second * 10
	c.Connection.SetDeadline(time.Now().Add(timeout))

	// send the Info request
	discoverMsg := []byte{3, 0, packetType}
	sendLen, err := c.Connection.Write(discoverMsg)
	if err != nil {
		return nil, err
	}
	log.Debug("Sent ", sendLen, " bytes")

	for {
		// Start a loop to read from UDP, and exit out if an error or timeout is experienced
		readLen, _, err = c.Connection.ReadFromUDP(inBuf)

		if err != nil {
			// Some kind of connection error
			return nil, err
		}

		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			// Timeout error
			return nil, err
		}

		if readLen != 0 {
			// Success!
			log.Debug("Read ", readLen, " bytes")
			log.Debug("Response hexadecimal: ", hex.EncodeToString(inBuf))
			break
		}
	}

	buffer := bytes.NewBuffer(inBuf)
	retLength := int(binary.LittleEndian.Uint16(buffer.Next(2)))
	if readLen != retLength {
		return nil, errors.New(fmt.Sprint("invalid reported buffer length: got ", retLength, ", expected ", readLen))
	}

	retType := uint8(buffer.Next(1)[0])
	if expect != retType {
		return nil, errors.New(fmt.Sprint("unexpected response packet type: got ", retType, ", expected ", expect))
	}

	return buffer, nil
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
// Connections time out after 10 seconds
func ScanServer(host string, port int) (serverstate OpenttdServerState, err error) {
	obj := OpenttdClientConnection{Hostname: host, Port: port}
	err = obj.Open()
	if err != nil {
		return OpenttdServerState{Status: false, Error: err}, err
	}

	// Let's get the initial set of data using CLIENT_FIND_SERVER
	result, err := obj.Query(openttd_packets_udp.ClientFindServer, openttd_packets_udp.ServerResponse)
	if err != nil {
		return OpenttdServerState{Status: false, Error: err}, err
	}
	serverstate = OpenttdServerState{Host: obj.Hostname}
	serverstate.populateServerState(result)

	// Then we get the company data using CLIENT_DETAIL_INFO
	result, err = obj.Query(openttd_packets_udp.ClientDetailInfo, openttd_packets_udp.ServerDetailInfo)
	if err != nil {
		return OpenttdServerState{Status: false, Error: err}, err
	}
	serverstate.populateCompanyState(result)
	return
}

// PopulateServerState populates an OpenttdServerState struct with data parsed from the Info request in the buffer.
func (server *OpenttdServerState) populateServerState(buf *bytes.Buffer) {
	// Props to https://github.com/vorot93/grokstat/blob/master/protocol_openttds.go for most of this code

	var protocolVer = int(buf.Next(1)[0])

	var activeNewGRFsNum int
	var activeNewGRFsInfo []OpenttdNewgrf
	if protocolVer >= 4 {
		activeNewGRFsNum = int(buf.Next(1)[0])
		for n := 0; n < activeNewGRFsNum; n++ {
			NewGRFID := getByteString(buf.Next(4))
			NewGRFMD5 := getByteString(buf.Next(16))
			activeNewGRFsInfo = append(activeNewGRFsInfo, OpenttdNewgrf{Identifier: NewGRFID, Hash: NewGRFMD5})
		}
	}

	var timeCurrent time.Time
	var timeStart time.Time
	if protocolVer >= 3 {
		timeCurrent = OttdDateFormat(binary.LittleEndian.Uint32(buf.Next(4)))
		timeStart = OttdDateFormat(binary.LittleEndian.Uint32(buf.Next(4)))
	}

	var maxCompanies *int
	var currentCompanies *int
	var maxSpectators *int
	if protocolVer >= 2 {
		maxCompanies = goutil.IntP(int(buf.Next(1)[0]))
		currentCompanies = goutil.IntP(int(buf.Next(1)[0]))
		maxSpectators = goutil.IntP(int(buf.Next(1)[0]))
	}
	serverNameBytes, _ := buf.ReadBytes(byte(0))
	serverName := string(bytes.Trim(serverNameBytes, "\x00"))

	serverVersionBytes, _ := buf.ReadBytes(byte(0))
	serverVersion := string(bytes.Trim(serverVersionBytes, "\x00"))

	languageId := int(buf.Next(1)[0])
	needPass := bool(buf.Next(1)[0] != 0)
	maxClients := int(buf.Next(1)[0])
	currentClients := int(buf.Next(1)[0])
	currentSpectators := int(buf.Next(1)[0])

	if protocolVer < 3 {
		_ = buf.Next(2)
		_ = buf.Next(2)
	}

	mapNameBytes, _ := buf.ReadBytes(byte(0))
	mapName := string(bytes.Trim(mapNameBytes, "\x00"))

	mapWidth := binary.LittleEndian.Uint16(buf.Next(2))
	mapHeight := binary.LittleEndian.Uint16(buf.Next(2))

	mapSet := int(buf.Next(1)[0])
	dedicatedServer := int(buf.Next(1)[0])

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

// PopulateCompanyState populates an OpenttdServerState struct with company data.
func (server *OpenttdServerState) populateCompanyState(buf *bytes.Buffer) {
	// Props to https://github.com/sonicsnes/node-gamedig/blob/master/protocols/openttd.js for most of this code

	var protocolVer = int(buf.Next(1)[0])
	if protocolVer >= 6 {
		var companies []OpenttdCompany
		numCompanies := int(buf.Next(1)[0])
		for i := 0; i < numCompanies; i++ {
			company := OpenttdCompany{}
			company.Id = uint8(buf.Next(1)[0])

			rawCompanyName, _ := buf.ReadBytes(byte(0))
			company.Name = string(bytes.Trim(rawCompanyName, "\x00"))

			company.YearStart = binary.LittleEndian.Uint32(buf.Next(4))
			company.Value = binary.LittleEndian.Uint64(buf.Next(8))
			company.Money = binary.LittleEndian.Uint64(buf.Next(8))
			company.Income = binary.LittleEndian.Uint64(buf.Next(8))
			company.Performance = binary.LittleEndian.Uint16(buf.Next(2))
			company.Passworded = !(int(buf.Next(1)[0]) == 0)

			company.Vehicles = OpenttdTypeCounts{}
			company.Vehicles.Train = binary.LittleEndian.Uint16(buf.Next(2))
			company.Vehicles.Truck = binary.LittleEndian.Uint16(buf.Next(2))
			company.Vehicles.Bus = binary.LittleEndian.Uint16(buf.Next(2))
			company.Vehicles.Aircraft = binary.LittleEndian.Uint16(buf.Next(2))
			company.Vehicles.Ship = binary.LittleEndian.Uint16(buf.Next(2))

			company.Stations = OpenttdTypeCounts{}
			company.Stations.Train = binary.LittleEndian.Uint16(buf.Next(2))
			company.Stations.Truck = binary.LittleEndian.Uint16(buf.Next(2))
			company.Stations.Bus = binary.LittleEndian.Uint16(buf.Next(2))
			company.Stations.Aircraft = binary.LittleEndian.Uint16(buf.Next(2))
			company.Stations.Ship = binary.LittleEndian.Uint16(buf.Next(2))

			clientsBytes, _ := buf.ReadBytes(byte(0))
			clients := bytes.Trim(clientsBytes, "\x00")
			log.Debug(clients)

			companies = append(companies, company)
		}
		server.Companies = companies
	} else {
		log.Warn("Unable to decode company details on Protocol Version ", protocolVer)
	}

}
