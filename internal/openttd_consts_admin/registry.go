package openttd_consts_admin

import (
	"errors"
	"fmt"
)

func GetAdminPacketType(p AdminPacket) (uint8, error) {
	switch p.(type) {
	case AdminJoin:
		return 0, nil
	case AdminQuit:
		return 1, nil
	}
	return 255, errors.New(fmt.Sprint("unknown packet type ", p))
}

const AdminUpdateFrequency = 2
const AdminPoll = 3
const AdminChat = 4
const AdminRcon = 5
const AdminGamescript = 6
const AdminPing = 7

// Server packets
const ServerFull = 100
const ServerBanned = 101
const ServerError = 102
const ServerNewgame = 105
const ServerShutdown = 106

// Server update packets
const ServerDate = 107
const ServerClientJoin = 108
const ServerClientInfo = 109
const ServerClientUpdate = 110
const ServerClientQuit = 111
const ServerClientError = 112
const ServerCompanyNew = 113
const ServerCompanyInfo = 114
const ServerCompanyUpdate = 115
const ServerCompanyRemove = 116
const ServerCompanyEconomy = 117
const ServerCompanyStats = 118
const ServerChat = 119
const ServerRcon = 120
const ServerConsole = 121
const ServerCmdNames = 122
const ServerCmdLogging = 123
const ServerGamescript = 124
const ServerRconEnd = 125
const ServerPong = 126
