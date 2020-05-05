package admin

// Client packets
const iAdminJoin uint8 = 0
const iAdminQuit uint8 = 1
const iAdminUpdateFrequency uint8 = 2
const iAdminPoll uint8 = 3
const iAdminChat uint8 = 4
const iAdminRcon uint8 = 5
const iAdminGamescript uint8 = 6
const iAdminPing uint8 = 7

// Server packets
const iServerFull uint8 = 100
const iServerBanned uint8 = 101
const iServerError uint8 = 102
const iServerProtocol uint8 = 103
const iServerWelcome uint8 = 104
const iServerNewgame uint8 = 105
const iServerShutdown uint8 = 106
const iServerDate uint8 = 107
const iServerClientJoin uint8 = 108
const iServerClientInfo uint8 = 109
const iServerClientUpdate uint8 = 110
const iServerClientQuit uint8 = 111
const iServerClientError uint8 = 112
const iServerCompanyNew uint8 = 113
const iServerCompanyInfo uint8 = 114
const iServerCompanyUpdate uint8 = 115
const iServerCompanyRemove uint8 = 116
const iServerCompanyEconomy uint8 = 117
const iServerCompanyStats uint8 = 118
const iServerChat uint8 = 119
const iServerRcon uint8 = 120
const iServerConsole uint8 = 121
const iServerCmdNames uint8 = 122
const iServerCmdLogging uint8 = 123
const iServerGamescript uint8 = 124
const iServerRconEnd uint8 = 125
const iServerPong uint8 = 126
