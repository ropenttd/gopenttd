package global

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ropenttd/gopenttd/internal/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/skybon/goutil"
	"time"
)

// populateServerState populates an OpenttdServerState struct with data parsed from the Info request in the buffer.
func (server *OpenttdServerState) PopulateServerState(buf *bytes.Buffer) {
	// Props to https://github.com/vorot93/grokstat/blob/master/protocol_openttds.go for most of this code

	var protocolVer = int(buf.Next(1)[0])

	var activeNewGRFsNum int
	var activeNewGRFsInfo []OpenttdNewgrf
	if protocolVer >= 4 {
		activeNewGRFsNum = int(buf.Next(1)[0])
		for n := 0; n < activeNewGRFsNum; n++ {
			NewGRFID := helpers.GetByteString(buf.Next(4))
			NewGRFMD5 := helpers.GetByteString(buf.Next(16))
			activeNewGRFsInfo = append(activeNewGRFsInfo, OpenttdNewgrf{Identifier: NewGRFID, Hash: NewGRFMD5})
		}
	}

	var timeCurrent time.Time
	var timeStart time.Time
	if protocolVer >= 3 {
		timeCurrent = helpers.OttdDateFormat(binary.LittleEndian.Uint32(buf.Next(4)))
		timeStart = helpers.OttdDateFormat(binary.LittleEndian.Uint32(buf.Next(4)))
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

// populateCompanyState populates an OpenttdServerState struct with company data.
func (server *OpenttdServerState) PopulateCompanyState(buf *bytes.Buffer) {
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
