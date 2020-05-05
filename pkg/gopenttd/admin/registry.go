package admin

import (
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
)

func GetRequestPacketType(p packets.AdminRequestPacket) (uint8, error) {
	switch p.(type) {
	case packets.AdminJoin:
		return 0, nil
	case packets.AdminQuit:
		return 1, nil
	case packets.AdminUpdateFrequency:
		return 2, nil
	case packets.AdminPoll:
		return 3, nil
	case packets.AdminChat:
		return 4, nil
	case packets.AdminRcon:
		return 5, nil
	case packets.AdminGamescript:
		return 6, nil
	case packets.AdminPing:
		return 7, nil
	}
	return 255, errors.New(fmt.Sprint("unknown packet type ", p))
}
