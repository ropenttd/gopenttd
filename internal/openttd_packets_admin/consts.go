package openttd_packets_admin

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

func (i AdminPacketIndex) IsRequest() bool {
	return i >= PacketAdminJoin && i <= PacketAdminPing
}
func (i AdminPacketIndex) IsResponse() bool {
	return i >= PacketServerFull && i <= PacketServerPong
}
