package gopenttd

import (
	"github.com/ropenttd/gopenttd/internal/helpers"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/consts"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
	log "github.com/sirupsen/logrus"
)

var state = OpenttdServerState{}

func handleServerWelcome(packet *packets.ServerWelcome, conn *OpenttdAdminConnection) {
	state.Name = packet.Name
	state.Version = packet.Version
	state.Dedicated = packet.Dedicated
	state.Map = packet.Map
	state.Seed = packet.Seed
	state.Environment = OpenttdEnvironment(packet.Landscape)
	state.DateStart = helpers.OttdDateFormat(packet.StartDate)
	state.MapWidth = packet.MapWidth
	state.MapHeight = packet.MapHeight
	return
}

func handleServerDate(packet *packets.ServerDate, conn *OpenttdAdminConnection) {
	state.DateCurrent = helpers.OttdDateFormat(packet.CurrentDate)
	return
}

func handleClientInfo(packet *packets.ServerClientInfo, conn *OpenttdAdminConnection) {
	return
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
func ScanServerAdm(host string, port int, password string) (err error) {
	obj := NewAdminConnection(host, port, password, "gopenttd")

	obj.RegWelcome(handleServerWelcome)
	obj.RegDate(handleServerDate)
	obj.RegClientInfo(handleClientInfo)
	err = obj.Open()
	if err != nil {
		return err
	}
	defer obj.Close()

	obj.Writer.Write(packets.AdminPoll{
		Type: consts.UpdateTypeDate,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: consts.UpdateTypeClientInfo,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: consts.UpdateTypeCompanyInfo,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: consts.UpdateTypeCompanyEconomy,
		ID:   ^uint32(0),
	})

	obj.Writer.Write(packets.AdminPoll{
		Type: consts.UpdateTypeCompanyStats,
		ID:   ^uint32(0),
	})

	for {
		packet, err := obj.ReadPacket()
		if err != nil {
			log.Error(err)
			break
		}
		log.Info(packet)
	}

	return
}
