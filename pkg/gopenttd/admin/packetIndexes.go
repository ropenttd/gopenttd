package admin

type AdminPacketIndex uint8

// Client packets
const PacketAdminJoin AdminPacketIndex = 0
const PacketAdminQuit AdminPacketIndex = 1
const PacketAdminUpdateFrequency AdminPacketIndex = 2
const PacketAdminPoll AdminPacketIndex = 3
const PacketAdminChat AdminPacketIndex = 4
const PacketAdminRcon AdminPacketIndex = 5
const PacketAdminGamescript AdminPacketIndex = 6
const PacketAdminPing AdminPacketIndex = 7

// Server packets
const PacketServerFull AdminPacketIndex = 100
const PacketServerBanned AdminPacketIndex = 101
const PacketServerError AdminPacketIndex = 102
const PacketServerProtocol AdminPacketIndex = 103
const PacketServerWelcome AdminPacketIndex = 104
const PacketServerNewgame AdminPacketIndex = 105
const PacketServerShutdown AdminPacketIndex = 106
const PacketServerDate AdminPacketIndex = 107
const PacketServerClientJoin AdminPacketIndex = 108
const PacketServerClientInfo AdminPacketIndex = 109
const PacketServerClientUpdate AdminPacketIndex = 110
const PacketServerClientQuit AdminPacketIndex = 111
const PacketServerClientError AdminPacketIndex = 112
const PacketServerCompanyNew AdminPacketIndex = 113
const PacketServerCompanyInfo AdminPacketIndex = 114
const PacketServerCompanyUpdate AdminPacketIndex = 115
const PacketServerCompanyRemove AdminPacketIndex = 116
const PacketServerCompanyEconomy AdminPacketIndex = 117
const PacketServerCompanyStats AdminPacketIndex = 118
const PacketServerChat AdminPacketIndex = 119
const PacketServerRcon AdminPacketIndex = 120
const PacketServerConsole AdminPacketIndex = 121
const PacketServerCmdNames AdminPacketIndex = 122
const PacketServerCmdLogging AdminPacketIndex = 123
const PacketServerGamescript AdminPacketIndex = 124
const PacketServerRconEnd AdminPacketIndex = 125
const PacketServerPong AdminPacketIndex = 126
