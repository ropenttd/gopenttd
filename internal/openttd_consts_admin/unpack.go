package openttd_consts_admin

import (
	"bytes"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

func Unpack(data []byte) (packet AdminPacket, err error) {
	reader := bytes.NewBuffer(data)
	packetType := uint8(reader.Next(1)[0])
	log.Debug(packetType)
	switch packetType {
	case 103:
		obj := ServerProtocol{Settings: map[uint16]uint16{}}
		binary.Read(reader, binary.LittleEndian, &obj.Version)
		var next bool
		binary.Read(reader, binary.LittleEndian, &next)
		for next {
			// there are settings to read
			var settingSlot uint16
			var settingVal uint16
			binary.Read(reader, binary.LittleEndian, &settingSlot)
			binary.Read(reader, binary.LittleEndian, &settingVal)
			obj.Settings[settingSlot] = settingVal
			binary.Read(reader, binary.LittleEndian, &next)
		}
		packet = obj
	case 104:
		obj := ServerWelcome{}
		serverNameBytes, _ := reader.ReadBytes(byte(0))
		obj.Name = string(bytes.Trim(serverNameBytes, "\x00"))

		serverVersionBytes, _ := reader.ReadBytes(byte(0))
		obj.Version = string(bytes.Trim(serverVersionBytes, "\x00"))

		binary.Read(reader, binary.LittleEndian, &obj.Dedicated)
		serverMapBytes, _ := reader.ReadBytes(byte(0))
		obj.Map = string(bytes.Trim(serverMapBytes, "\x00"))

		binary.Read(reader, binary.LittleEndian, &obj.Seed)
		binary.Read(reader, binary.LittleEndian, &obj.Landscape)
		binary.Read(reader, binary.LittleEndian, &obj.StartDate)
		binary.Read(reader, binary.LittleEndian, &obj.MapWidth)
		binary.Read(reader, binary.LittleEndian, &obj.MapHeight)

		packet = obj
	}
	return packet, nil
}
