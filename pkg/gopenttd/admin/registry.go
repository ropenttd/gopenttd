package admin

import (
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
)

func GetRequestPacketType(p packets.AdminRequestPacket) (uint8, error) {
	switch p.(type) {
	case packets.AdminJoin:
		return iAdminJoin, nil
	case packets.AdminQuit:
		return iAdminQuit, nil
	case packets.AdminUpdateFrequency:
		return iAdminUpdateFrequency, nil
	case packets.AdminPoll:
		return iAdminPoll, nil
	case packets.AdminChat:
		return iAdminChat, nil
	case packets.AdminRcon:
		return iAdminRcon, nil
	case packets.AdminGamescript:
		return iAdminGamescript, nil
	case packets.AdminPing:
		return iAdminPing, nil
	}
	return 255, errors.New(fmt.Sprint("unknown packet type ", p))
}
