package openttd_packets_udp

// This is a list of constants defining the packet order for each of the given packets.
// Source: https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/udp.h

type UdpPacketIndex uint8

const (
	// Client-server related
	ClientFindServer UdpPacketIndex = iota
	ServerResponse
	ClientDetailInfo
	ServerDetailInfo

	// Master server related
	ServerRegister
	MasterAckRegister
	ClientGetList
	MasterResponseList
	ServerUnregister

	// Client-server related
	ClientGetNewgrfs
	ServerNewgrfs

	// Master server related
	MasterSessionkey

	end
)
