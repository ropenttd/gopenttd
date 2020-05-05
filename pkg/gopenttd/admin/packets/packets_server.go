package packets

import (
	"bytes"
	"encoding/binary"
)

type AdminResponsePacket interface {
}

// As defined in https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_admin.h

type ServerFull struct { // Type 100
}

type ServerBanned struct { // Type 101

}

type ServerError struct { // Type 102
	ErrorCode uint8 // NetworkErrorCode the error caused.
}

type ServerProtocol struct { // Type 103
	Version  uint8
	Settings map[uint16]uint16
}

func (p *ServerProtocol) Unpack(buffer *bytes.Buffer) (err error) {
	var ver uint8
	binary.Read(buffer, binary.LittleEndian, &ver)
	p.Version = ver
	var next bool
	binary.Read(buffer, binary.LittleEndian, &next)
	for next {
		// there are settings to read
		var settingSlot uint16
		var settingVal uint16
		binary.Read(buffer, binary.LittleEndian, &settingSlot)
		binary.Read(buffer, binary.LittleEndian, &settingVal)
		p.Settings[settingSlot] = settingVal
		binary.Read(buffer, binary.LittleEndian, &next)
	}
	return nil
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

type ServerNewgame struct { // Type 105

}

type ServerShutdown struct { // Type 106

}

// Update packets
type ServerDate struct { // Type 107
	StartDate uint32 // Current game date.
}

type ServerClientJoin struct { // Type 108

	ID uint32 // ID of the new client.
}

type ServerClientInfo struct { // Type 109

	ID       uint32 // ID of the client.
	Address  string // Network address of the client.
	Name     string // Name of the client.
	Language uint8  // Language of the client.
	JoinDate uint32 // Date the client joined the game.
	Company  uint8  // ID of the company the client is playing as (255 for spectators).
}

type ServerClientUpdate struct { // Type 110

	ID      uint32 // ID of the client.
	Name    string // Name of the client.
	Company uint8  // ID of the company the client is playing as (255 for spectators).
}

type ServerClientQuit struct { // Type 111

	ID uint32 // ID of the leaving client.
}

type ServerClientError struct { // Type 112

	ID    uint32 // ID of the client throwing the error.
	Error uint8  // Error the client made (see NetworkErrorCode).
}

type ServerCompanyNew struct { // Type 113

	ID uint8 // ID of the new company.
}

type ServerCompanyInfo struct { // Type 114

	ID        uint8  // ID of the company.
	Name      string // Name of the company.
	Manager   string // Name of the companies manager.
	Colour    uint8  // Main company colour.
	Password  bool   // Company is password protected.
	StartDate uint32 // Year the company was inaugurated.
	IsAI      bool   // Company is an AI.
}

type ServerCompanyUpdate struct { // Type 115

	ID                 uint8  // ID of the company.
	Name               string // Name of the company.
	Manager            string // Name of the companies manager.
	Colour             uint8  // Main company colour.
	Password           bool   // Company is password protected.
	BankruptcyQuarters uint8  // Quarters of Bankruptcy.
	Share1             uint8  // Owner of Share 1.
	Share2             uint8  // Owner of Share 2.
	Share3             uint8  // Owner of Share 3.
	Share4             uint8  // Owner of Share 4.
}

type ServerCompanyRemove struct { // Type 116

	ID     uint8 // ID of the company.
	Reason uint8 // Reason for being removed (see #AdminCompanyRemoveReason).
}

type ServerCompanyEconomy struct { // Type 117

	ID                         uint8  // ID of the company.
	Money                      uint64 // Money (cash in hand).
	Loan                       uint64 // Loan.
	Income                     int64  // Income.
	CargoThisQuarter           uint16 // Delivered cargo (this quarter).
	ValueLastQuarter           uint64 // Company value (last quarter).
	PerformanceLastQuarter     uint16 // Performance (last quarter).
	CargoLastQuarter           uint16 // Delivered cargo (last quarter).
	ValuePreviousQuarter       uint64 // Company value (previous quarter).
	PerformancePreviousQuarter uint16 // Performance (previous quarter).
	CargoPreviousQuarter       uint16 // Delivered cargo (previous quarter).
}

type ServerCompanyStats struct { // Type 118

	ID            uint8  // ID of the company.
	Trains        uint16 // Number of trains.
	Lorries       uint16 // Number of lorries.
	Buses         uint16 // Number of busses.
	Planes        uint16 // Number of planes.
	Ships         uint16 // Number of ships.
	TrainStations uint16 // Number of train stations.
	LorryStations uint16 // Number of lorry stations.
	BusStops      uint16 // Number of bus stops.
	Airports      uint16 // Number of airports and heliports.
	Harbours      uint16 // Number of harbours.
}

type ServerChat struct { // Type 119

	Action      uint8  // Action such as NETWORK_ACTION_CHAT_CLIENT (see #NetworkAction).
	Destination uint8  // Destination type such as DESTTYPE_BROADCAST (see #DestType).
	ID          uint32 // ID of the client who sent this message.
	Message     string // Message.
	Money       uint64 // Money (only when it is a 'give money' action).
}

type ServerRcon struct { // Type 120

	Colour uint8  // Colour as it would be used on the server or a client.
	Output string // Output of the executed command.
}

type ServerConsole struct { // Type 121

	Origin  string // The origin of the text, e.g. "console" for console, or "net" for network related (debug) messages.
	Message string // Text as found on the console of the server.
}

type ServerCmdNames struct { // Type 122
	/*
	* NOTICE: Pack provided with this packet is not stable and will not be
	*         treated as such. Do not rely on IDs or names to be constant
	*         across different versions / revisions of OpenTTD.
	*         Pack provided in this packet is for logging purposes only.
	 */

	More bool   // Pack to follow.
	ID   uint16 // ID of the DoCommand.
	Name string // Name of the DoCommand.
}

type ServerCmdLogging struct { // Type 123
	/*
	* NOTICE: Pack provided with this packet is not stable and will not be
	*         treated as such. Do not rely on IDs or names to be constant
	*         across different versions / revisions of OpenTTD.
	*         Pack provided in this packet is for logging purposes only.
	 */

	Client    uint32 // ID of the client sending the command.
	Company   uint8  // ID of the company (0..MAX_COMPANIES-1).
	CommandID uint16 // ID of the command.
	V1        uint32 // P1 (variable data passed to the command).
	V2        uint32 // P2 (variable data passed to the command).
	Tile      uint32 // Tile where this is taking place.
	Message   string // Text passed to the command.
	Frame     uint32 // Frame of execution.
}

type ServerGamescript struct { // Type 124
	// Pack on this isn't available in the source? Making an assumption.

	Json string // JSON string from the GameScript.
}

type ServerRconEnd struct { // Type 125

	Command string // The command as requested by the admin connection.
}

type ServerPong struct { // Type 126

	Token uint32 // Integer value requested in the Ping.
}
