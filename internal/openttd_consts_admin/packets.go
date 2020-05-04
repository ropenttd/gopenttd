package openttd_consts_admin

import "bytes"

// As defined in https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_admin.h

type AdminPacket interface {
	Data() bytes.Buffer
}

// Client packets
type AdminJoin struct { // Type 0
	Password   string
	ClientName string
	Version    string
}

func (p AdminJoin) Data() (out bytes.Buffer) {
	out.WriteString(p.Password)
	out.WriteByte(byte(0))
	out.WriteString(p.ClientName)
	out.WriteByte(byte(0))
	out.WriteString(p.Version)
	out.WriteByte(byte(0))
	return out
}

type AdminQuit struct { // Type 1
}

func (p AdminQuit) Data() (out bytes.Buffer) {
	return out
}

type ServerProtocol struct { // Type 103
	Version  uint8
	Settings map[uint16]uint16
}

func (p ServerProtocol) Data() (out bytes.Buffer) {
	// this is an incoming packet, we don't serialize it back
	return out
}

type ServerWelcome struct { // Type 104
	Name      string // Name of the Server (e.g. as advertised to master server).
	Version   string // OpenTTD version string.
	Dedicated bool   // Server is dedicated.
	Map       string // Name of the Map.
	Seed      uint32 // Random seed of the Map.
	Landscape uint8  // Landscape of the Map.
	StartDate uint32 // Start date of the Map.
	MapWidth  uint16 // Map width.
	MapHeight uint16 // Map height.
}

func (p ServerWelcome) Data() (out bytes.Buffer) {
	// this is an incoming packet, we can't serialize it back
	return out
}
