package openttd_packets_admin

import (
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/admin/packets"
)

type AdminPacketIndex uint8

const (
	// Client packets
	PacketAdminJoin AdminPacketIndex = iota
	PacketAdminQuit
	PacketAdminUpdateFrequency
	PacketAdminPoll
	PacketAdminChat
	PacketAdminRcon
	PacketAdminGamescript
	PacketAdminPing
)

const (
	// Server packets
	PacketServerFull AdminPacketIndex = 100 + iota
	PacketServerBanned
	PacketServerError
	PacketServerProtocol
	PacketServerWelcome
	PacketServerNewgame
	PacketServerShutdown
	PacketServerDate
	PacketServerClientJoin
	PacketServerClientInfo
	PacketServerClientUpdate
	PacketServerClientQuit
	PacketServerClientError
	PacketServerCompanyNew
	PacketServerCompanyInfo
	PacketServerCompanyUpdate
	PacketServerCompanyRemove
	PacketServerCompanyEconomy
	PacketServerCompanyStats
	PacketServerChat
	PacketServerRcon
	PacketServerConsole
	PacketServerCmdNames
	PacketServerCmdLogging
	PacketServerGamescript
	PacketServerRconEnd
	PacketServerPong
)

func GetRequestPacketType(p packets.AdminRequestPacket) (AdminPacketIndex, error) {
	switch p.(type) {
	case packets.AdminJoin:
		return PacketAdminJoin, nil
	case packets.AdminQuit:
		return PacketAdminQuit, nil
	case packets.AdminUpdateFrequency:
		return PacketAdminUpdateFrequency, nil
	case packets.AdminPoll:
		return PacketAdminPoll, nil
	case packets.AdminChat:
		return PacketAdminChat, nil
	case packets.AdminRcon:
		return PacketAdminRcon, nil
	case packets.AdminGamescript:
		return PacketAdminGamescript, nil
	case packets.AdminPing:
		return PacketAdminPing, nil
	}
	return 255, errors.New(fmt.Sprint("unknown packet type ", p))
}
