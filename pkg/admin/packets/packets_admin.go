package packets

import (
	"bytes"
	"encoding/binary"
	"github.com/ropenttd/gopenttd/internal/helpers"
	"github.com/ropenttd/gopenttd/internal/openttd_packets_admin"
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
)

// As defined in https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_admin.h

type AdminRequestPacket interface {
	Pack() bytes.Buffer
	// Unpack(*bytes.Buffer) error
	PacketType() openttd_packets_admin.AdminPacketIndex
}

// Client packets
type AdminJoin struct { // Type 0
	Password   string
	ClientName string
	Version    string
}

func (p AdminJoin) Pack() (out bytes.Buffer) {
	out.Write(helpers.PackString(p.Password))
	out.Write(helpers.PackString(p.ClientName))
	out.Write(helpers.PackString(p.Version))
	return out
}
func (p AdminJoin) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminJoin
}

type AdminQuit struct { // Type 1
}

func (p AdminQuit) Pack() (out bytes.Buffer) {
	return out
}
func (p AdminQuit) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminQuit
}

type AdminUpdateFrequency struct { // Type 2
	Type      enum.UpdateType      // Update type (see #AdminUpdateType).
	Frequency enum.UpdateFrequency // Update frequency (see #AdminUpdateFrequency), setting #ADMIN_FREQUENCY_POLL is always ignored.
}

func (p AdminUpdateFrequency) Pack() (out bytes.Buffer) {
	binary.Write(&out, binary.LittleEndian, p.Type)
	binary.Write(&out, binary.LittleEndian, p.Frequency)
	return out
}
func (p AdminUpdateFrequency) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminUpdateFrequency
}

type AdminPoll struct { // Type 3
	Type enum.UpdateType // #AdminUpdateType the server should answer for, only if #AdminUpdateFrequency #ADMIN_FREQUENCY_POLL is advertised in the PROTOCOL packet.
	ID   uint32          // uint32  ID relevant to the packet type, e.g.
	// - the client ID for #ADMIN_UPDATE_CLIENT_INFO. Use UINT32_MAX to show all clients.
	// - the company ID for #ADMIN_UPDATE_COMPANY_INFO. Use UINT32_MAX to show all companies.
}

func (p AdminPoll) Pack() (out bytes.Buffer) {
	binary.Write(&out, binary.LittleEndian, p.Type)
	binary.Write(&out, binary.LittleEndian, p.ID)
	return out
}
func (p AdminPoll) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminPoll
}

type AdminChat struct { // Type 4
	Action        enum.Action      // Action such as NETWORK_ACTION_CHAT_CLIENT (see #NetworkAction).
	Destination   enum.Destination // Destination type such as DESTTYPE_BROADCAST (see #DestType).
	DestinationID uint32           // ID of the destination such as company or client id.
	Message       string           // Message.
}

func (p AdminChat) Pack() (out bytes.Buffer) {
	binary.Write(&out, binary.LittleEndian, p.Action)
	binary.Write(&out, binary.LittleEndian, p.Destination)
	binary.Write(&out, binary.LittleEndian, p.DestinationID)
	out.Write(helpers.PackString(p.Message))
	return out
}
func (p AdminChat) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminChat
}

type AdminRcon struct { // Type 5
	Command string // Command to be executed.
}

func (p AdminRcon) Pack() (out bytes.Buffer) {
	out.Write(helpers.PackString(p.Command))
	return out
}
func (p AdminRcon) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminRcon
}

type AdminGamescript struct { // Type 6
	Json string // JSON string for the GameScript.
}

func (p AdminGamescript) Pack() (out bytes.Buffer) {
	out.Write(helpers.PackString(p.Json))
	return out
}
func (p AdminGamescript) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminGamescript
}

type AdminPing struct { // Type 7
	Token uint32 // Integer value to pass to the server, which is quoted in the reply.
}

func (p AdminPing) Pack() (out bytes.Buffer) {
	binary.Write(&out, binary.LittleEndian, p.Token)
	return out
}
func (p AdminPing) PacketType() openttd_packets_admin.AdminPacketIndex {
	return openttd_packets_admin.PacketAdminPing
}
