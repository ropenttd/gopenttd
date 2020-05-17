package admin

import (
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
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
