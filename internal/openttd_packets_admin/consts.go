package openttd_packets_admin

import (
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/admin/packets"
)

type adminPacketIndex uint8

const (
	// Client packets
	packetAdminJoin adminPacketIndex = iota
	packetAdminQuit
	packetAdminUpdateFrequency
	packetAdminPoll
	packetAdminChat
	packetAdminRcon
	packetAdminGamescript
	packetAdminPing
)

const (
	// Server packets
	packetServerFull adminPacketIndex = 100 + iota
	packetServerBanned
	packetServerError
	packetServerProtocol
	packetServerWelcome
	packetServerNewgame
	packetServerShutdown
	packetServerDate
	packetServerClientJoin
	packetServerClientInfo
	packetServerClientUpdate
	packetServerClientQuit
	packetServerClientError
	packetServerCompanyNew
	packetServerCompanyInfo
	packetServerCompanyUpdate
	packetServerCompanyRemove
	packetServerCompanyEconomy
	packetServerCompanyStats
	packetServerChat
	packetServerRcon
	packetServerConsole
	packetServerCmdNames
	packetServerCmdLogging
	packetServerGamescript
	packetServerRconEnd
	packetServerPong
)

func GetRequestPacketType(p packets.AdminRequestPacket) (adminPacketIndex, error) {
	switch p.(type) {
	case packets.AdminJoin:
		return packetAdminJoin, nil
	case packets.AdminQuit:
		return packetAdminQuit, nil
	case packets.AdminUpdateFrequency:
		return packetAdminUpdateFrequency, nil
	case packets.AdminPoll:
		return packetAdminPoll, nil
	case packets.AdminChat:
		return packetAdminChat, nil
	case packets.AdminRcon:
		return packetAdminRcon, nil
	case packets.AdminGamescript:
		return packetAdminGamescript, nil
	case packets.AdminPing:
		return packetAdminPing, nil
	}
	return 255, errors.New(fmt.Sprint("unknown packet type ", p))
}
