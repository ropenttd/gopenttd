package openttd_packets_udp

// This is a list of constants defining the packet order for each of the given packets.
// Source: https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/udp.h

// Client-server related
const ClientFindServer = 0
const ServerResponse = 1
const ClientDetailInfo = 2
const ServerDetailInfo = 3

// Master server related
const ServerRegister = 4
const MasterAckRegister = 5
const ClientGetList = 6
const MasterResponseList = 7
const ServerUnregister = 8

// Client-server related
const ClientGetNewgrfs = 9
const ServerNewgrfs = 10

// Master server related
const MasterSessionkey = 11

const end = 12
