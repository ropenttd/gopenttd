package openttd_packets_tcp

// This is a list of constants defining the packet order for each of the given packets.
// We don't use this in this library at present, but you might find it useful if you want to pretend to be a client for some reason.

// Source: https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_game.h

type TcpPacketIndex uint8

const (
	ServerFull TcpPacketIndex = iota
	ServerBanned
	ClientJoin
	ServerError
	ClientCompanyInfo
	ServerCompanyInfo
	ServerClientInfo
	ServerNeedPassword
	ClientPassword
	ServerWelcome
	ClientGetmap
	ServerWait
	ServerMap
	ClientMapOk
	ServerJoin
	ServerFrame
	ServerSync
	ClientAck
	ClientCommand
	ServerCommand
	ClientChat
	ServerChat
	ClientSetPassword
	ClientSetName
	ClientQuit
	ClientError
	ServerQuit
	ServerErrorQuit
	ServerShutdown
	ServerNewgame
	ServerRcon
	ClientRcon
	ServerCheckNewgrfs
	ClientNewgrfsChecked
	ServerMove
	ClientMove
	ServerCompanyUpdate
	ServerConfigEnd

	end
)
