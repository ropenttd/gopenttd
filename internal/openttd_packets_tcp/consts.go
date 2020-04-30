package openttd_packets_udp

// This is a list of constants defining the packet order for each of the given packets.
// We don't use this in this library at present, but you might find it useful if you want to pretend to be a client for some reason.

// Source: https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_game.h

const ServerFull = 0
const ServerBanned = 1
const ClientJoin = 2
const ServerError = 3
const ClientCompanyInfo = 4
const ServerCompanyInfo = 5
const ServerClientInfo = 6
const ServerNeedPassword = 7
const ClientPassword = 8
const ServerWelcome = 9
const ClientGetmap = 10
const ServerWait = 11
const ServerMap = 12
const ClientMapOk = 13
const ServerJoin = 14
const ServerFrame = 15
const ServerSync = 16
const ClientAck = 17
const ClientCommand = 18
const ServerCommand = 19
const ClientChat = 20
const ServerChat = 21
const ClientSetPassword = 22
const ClientSetName = 23
const ClientQuit = 24
const ClientError = 25
const ServerQuit = 26
const ServerErrorQuit = 27
const ServerShutdown = 28
const ServerNewgame = 29
const ServerRcon = 30
const ClientRcon = 31
const ServerCheckNewgrfs = 32
const ClientNewgrfsChecked = 33
const ServerMove = 34
const ClientMove = 35
const ServerCompanyUpdate = 36
const ServerConfigEnd = 37

const end = 38
