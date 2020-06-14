package admin

import (
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
)

// This file contains all the possible structs that can be
// handled by AddHandler/EventHandler.
// DO NOT ADD ANYTHING BUT EVENT HANDLER STRUCTS TO THIS FILE.
//go:generate go run tools/cmd/eventhandlers/main.go

// Connect is the data for a Connect event.
// This is a synthetic event and is not dispatched by OpenTTD.
type Connect struct{}

// Disconnect is the data for a Disconnect event.
// This is a synthetic event and is not dispatched by OpenTTD.
type Disconnect struct{}

// Event provides a basic initial struct for all game events.
type Event struct {
	Type    uint8  `json:"t"`
	RawData []byte `json:"d"`
	// Struct contains one of the other types in this file.
	Struct interface{} `json:"-"`
}

type Full struct { // Type 100
}
type Banned struct { // Type 101
}
type Error struct { // Type 102
	ErrorCode enum.NetError // NetworkErrorCode the error caused.
}

// Protocol fires when receiving a Protocol packet during join.
type Protocol struct { // Type 103
	Version  uint8
	Settings map[uint16]uint16
}

// Welcome fires when receiving a Welcome packet, usually during join or after map load.
type Welcome struct { // Type 104
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

type Newgame struct { // Type 105
}

type Shutdown struct { // Type 106
}

// Update packets

// Date fires when the game provides an update to the current game date.
type Date struct { // Type 107
	CurrentDate uint32 // Current game date.
}

// ClientJoin fires when a client joins.
type ClientJoin struct { // Type 108
	ID uint32 // ID of the new client.
}

// ClientInfo fires when information is provided for a client (usually polled).
type ClientInfo struct { // Type 109
	ID       uint32 // ID of the client.
	Address  string // Network address of the client. (Can be nil if they're the host)
	Name     string // Name of the client.
	Language uint8  // Language of the client.
	JoinDate uint32 // Date the client joined the game. (Can be nil if they're the host)
	Company  uint8  // ID of the company the client is playing as (255 for spectators).
}

// ClientUpdate fires when a client changes the company it is playing as, or its name.
type ClientUpdate struct { // Type 110
	ID      uint32 // ID of the client.
	Name    string // Name of the client.
	Company uint8  // ID of the company the client is playing as (255 for spectators).
}

// ClientQuit fires when a client leaves the game.
type ClientQuit struct { // Type 111
	ID uint32 // ID of the leaving client.
}

// ClientError fires when a client experiences an error (usually causes them to quit).
type ClientError struct { // Type 112
	ID    uint32        // ID of the client throwing the error.
	Error enum.NetError // Error the client made (see NetworkErrorCode).
}

// CompanyNew fires when a new company is founded.
type CompanyNew struct { // Type 113
	ID uint8 // ID of the new company.
}

// CompanyInfo fires when information for a company is provided (usually polled).
type CompanyInfo struct { // Type 114
	ID        uint8  // ID of the company.
	Name      string // Name of the company.
	Manager   string // Name of the companies manager.
	Colour    uint8  // Main company colour.
	Password  bool   // Company is password protected.
	StartDate uint32 // Year the company was inaugurated.
	IsAI      bool   // Company is an AI.
}

// CompanyUpdate fires when a company changes something about its metadata, such as its name.
type CompanyUpdate struct { // Type 115
	ID                 uint8  // ID of the company.
	Name               string // Name of the company.
	Manager            string // Name of the company's manager.
	Colour             uint8  // Main company colour.
	Password           bool   // Company is password protected.
	BankruptcyQuarters uint8  // Quarters of Bankruptcy.
	Share1             uint8  // Owner of Share 1.
	Share2             uint8  // Owner of Share 2.
	Share3             uint8  // Owner of Share 3.
	Share4             uint8  // Owner of Share 4.
}

// CompanyRemove fires when a company is removed from the game.
type CompanyRemove struct { // Type 116
	ID     uint8                    // ID of the company.
	Reason enum.CompanyRemoveReason // Reason for being removed.
}

// CompanyEconomy fires when new economical data for the company is available (either polled or regularly, depends what you signed up for)
type CompanyEconomy struct { // Type 117
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

// CompanyStats fires when new statistics for the company are available (either polled or regularly, depends what you signed up for)
type CompanyStats struct { // Type 118
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

// Chat fires when any new chat message is posted.
type Chat struct { // Type 119
	Action      uint8  // Action such as NETWORK_ACTION_CHAT_CLIENT (see #NetworkAction).
	Destination uint8  // Destination type such as DESTTYPE_BROADCAST (see #DestType).
	ID          uint32 // ID of the client who sent this message.
	Message     string // Message.
	Money       uint64 // Money (only when it is a 'give money' action).
}

// Rcon fires when a line of RCON output from the server is returned.
// Use in conjunction with RconEnd to determine when a command has finished.
type Rcon struct { // Type 120
	Colour uint8  // Colour as it would be used on the server or a client.
	Output string // Output of the executed command.
}

// Console fires when a console message is printed to the server output.
type Console struct { // Type 121
	Origin  string // The origin of the text, e.g. "console" for console, or "net" for network related (debug) messages.
	Message string // Text as found on the console of the server.
}

// CmdNames is a response to a request for command names.
// You probably shouldn't track this event - use the state on the session.
type CmdNames struct { // Type 122
	/*
	* NOTICE: Pack provided with this packet is not stable and will not be
	*         treated as such. Do not rely on IDs or names to be constant
	*         across different versions / revisions of OpenTTD.
	*         Pack provided in this packet is for logging purposes only.
	 */
	Commands map[uint16]string // Map of the ID of the DoCommand with the name of it.
}

// CmdLogging is an entry of some kind of game event that occurred. Very useful for auditing.
// You usually need realtime updates to get these.
type CmdLogging struct { // Type 123
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

// Gamescript is some data that was sent by a GameScript running on the server.
type Gamescript struct { // Type 124
	// Pack on this isn't available in the source? Making an assumption.
	Json string // JSON string from the GameScript.
}

// RconEnd denotes that the given command finished.
type RconEnd struct { // Type 125
	Command string // The command as requested by the admin connection.
}

// Pong is sent in response to a Ping request.
// Probably no need to track this one; gopenttd will handle ping/pong and timeouts for you.
type Pong struct { // Type 126
	Token uint32 // Integer value requested in the Ping.
}
