package admin

import (
	"github.com/ropenttd/gopenttd/internal/helpers"
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
	"github.com/ropenttd/gopenttd/pkg/admin/packets"
	"github.com/ropenttd/gopenttd/pkg/util"
	log "github.com/sirupsen/logrus"
)

// yes this is a globalvar sue me, it's the easiest way to provide this helper function
// todo: move this out to cmd, and find a better carrier for more data than OpenttdServerState (maybe OpenttdAdminServerState?)
var lstate *util.OpenttdServerState

func handleServerWelcome(packet *packets.ServerWelcome, conn *OpenttdAdminConnection) {
	lstate.Name = packet.Name
	lstate.Version = packet.Version
	lstate.Dedicated = packet.Dedicated
	lstate.Map = packet.Map
	lstate.Seed = packet.Seed
	lstate.Environment = util.OpenttdEnvironment(packet.Landscape)
	lstate.DateStart = helpers.OttdDateFormat(packet.StartDate)
	lstate.MapWidth = packet.MapWidth
	lstate.MapHeight = packet.MapHeight
	return
}

func handleServerDate(packet *packets.ServerDate, conn *OpenttdAdminConnection) {
	lstate.DateCurrent = helpers.OttdDateFormat(packet.CurrentDate)
	return
}

func handleClientInfo(packet *packets.ServerClientInfo, conn *OpenttdAdminConnection) {
	if packet.ID == 1 {
		// it's the server
		return
	}
	lstate.NumClients += 1
	if packet.Company == 255 {
		lstate.NumSpectators += 1
	}
	return
}

func handleCompanyInfo(packet *packets.ServerCompanyInfo, conn *OpenttdAdminConnection) {
	if lstate.Companies == nil {
		lstate.Companies = map[uint8]util.OpenttdCompany{}
	}
	company := util.OpenttdCompany{}
	company.Name = packet.Name
	company.YearStart = packet.StartDate
	company.Passworded = packet.Password
	lstate.Companies[packet.ID] = company
	return
}

func handleCompanyEconomy(packet *packets.ServerCompanyEconomy, conn *OpenttdAdminConnection) {
	company := lstate.Companies[packet.ID]
	company.Money = packet.Money
	company.Income = packet.Income
	company.Value = packet.ValueLastQuarter
	company.Performance = packet.PerformanceLastQuarter
	lstate.Companies[packet.ID] = company
	return
}

func handleCompanyStats(packet *packets.ServerCompanyStats, conn *OpenttdAdminConnection) {
	company := lstate.Companies[packet.ID]
	company.Vehicles.Train = packet.Trains
	company.Vehicles.Truck = packet.Lorries
	company.Vehicles.Bus = packet.Buses
	company.Vehicles.Aircraft = packet.Planes
	company.Vehicles.Ship = packet.Ships
	company.Stations.Train = packet.TrainStations
	company.Stations.Truck = packet.LorryStations
	company.Stations.Bus = packet.BusStops
	company.Stations.Aircraft = packet.Airports
	company.Stations.Ship = packet.Harbours
	lstate.Companies[packet.ID] = company
	return
}

// ScanServer takes a hostname and port and returns an OpenttdServerState.struct containing the data available from it.
func ScanServerAdm(host string, port int, password string) (state util.OpenttdServerState, err error) {
	obj := NewAdminConnection(host, port, password, "gopenttd")

	obj.RegWelcome(handleServerWelcome)
	obj.RegDate(handleServerDate)
	obj.RegClientInfo(handleClientInfo)
	obj.RegCompanyInfo(handleCompanyInfo)
	obj.RegCompanyEconomy(handleCompanyEconomy)
	obj.RegCompanyStats(handleCompanyStats)
	err = obj.Open()
	if err != nil {
		return util.OpenttdServerState{}, err
	}
	defer obj.Close()

	obj.Writer.Write(packets.AdminPoll{
		Type: enum.UpdateTypeDate,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: enum.UpdateTypeClientInfo,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: enum.UpdateTypeCompanyInfo,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: enum.UpdateTypeCompanyEconomy,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: enum.UpdateTypeCompanyStats,
		ID:   ^uint32(0),
	})

	// send a ping to denote the end of the request (hopefully we get the pong back after all the other data)
	// read also - this is hacky as all hell, please PR a saner way of doing this
	obj.Writer.Write(packets.AdminPing{
		Token: 65535,
	})

	lstate = &util.OpenttdServerState{}

	var cont bool = true
	for cont {
		packet, err := obj.ReadPacket()
		if err != nil {
			log.Error(err)
			break
		}
		log.Info(packet)
		switch packet.(type) {
		case *packets.ServerPong:
			cont = false
		}
	}

	return *lstate, nil
}
