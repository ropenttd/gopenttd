package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

type AdminResponsePacket interface {
	String() string
	Unmarshal(buffer *bytes.Buffer) error
}

func genericUnmarshal(p AdminResponsePacket, buffer *bytes.Buffer) (err error) {
	// Some funky bullshit that iterates through the elements in the ResponsePacket and unpacks them from the buffer
	// based on their type
	// I'm not even going to pretend to know how this works
	v := reflect.ValueOf(p).Elem()
	if !v.IsValid() {
		return errors.New(fmt.Sprint(v, "is invalid"))
	}
	for i := 0; i < v.NumField(); i++ {
		switch f := v.Field(i); f.Kind() {
		// binary.Read() doesn't appear to work here (always returns 0?) so do things the long way
		// i.e binary.Read(buffer, binary.LittleEndian, nv)
		case reflect.Bool:
			var nv bool
			nv = uint8(buffer.Next(1)[0]) != 0
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint8:
			var nv uint8
			nv = uint8(buffer.Next(1)[0])
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint16:
			var nv uint16
			nv = binary.LittleEndian.Uint16(buffer.Next(2))
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint32:
			var nv uint32
			nv = binary.LittleEndian.Uint32(buffer.Next(4))
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint64:
			var nv uint64
			nv = binary.LittleEndian.Uint64(buffer.Next(8))
			f.Set(reflect.ValueOf(nv))
		case reflect.String:
			nvBytes, _ := buffer.ReadBytes(byte(0))
			nv := string(bytes.Trim(nvBytes, "\x00"))
			f.Set(reflect.ValueOf(nv))
		}
	}
	return err
}

// As defined in https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_admin.h

type ServerFull struct { // Type 100
}

func (p ServerFull) String() string {
	return "SERVER_FULL"
}
func (p *ServerFull) Unmarshal(buffer *bytes.Buffer) (err error) {
	// no unpacking
	return nil
}

type ServerBanned struct { // Type 101
}

func (p ServerBanned) String() string {
	return "SERVER_BANNED"
}
func (p *ServerBanned) Unmarshal(buffer *bytes.Buffer) (err error) {
	// no unpacking
	return nil
}

type ServerError struct { // Type 102
	ErrorCode uint8 // NetworkErrorCode the error caused.
}

func (p ServerError) String() string {
	return "SERVER_ERROR"
}
func (p *ServerError) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerProtocol struct { // Type 103
	Version  uint8
	Settings map[uint16]uint16
}

func (p ServerProtocol) String() string {
	return "SERVER_PROTOCOL"
}
func (p *ServerProtocol) Unmarshal(buffer *bytes.Buffer) (err error) {
	var ver uint8
	binary.Read(buffer, binary.LittleEndian, &ver)
	p.Version = ver
	p.Settings = map[uint16]uint16{}
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

func (p ServerWelcome) String() string {
	return "SERVER_WELCOME"
}
func (p *ServerWelcome) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerNewgame struct { // Type 105
}

func (p ServerNewgame) String() string {
	return "SERVER_NEWGAME"
}
func (p *ServerNewgame) Unmarshal(buffer *bytes.Buffer) (err error) {
	// no unmarshalling
	return
}

type ServerShutdown struct { // Type 106
}

func (p ServerShutdown) String() string {
	return "SERVER_SHUTDOWN"
}
func (p *ServerShutdown) Unmarshal(buffer *bytes.Buffer) (err error) {
	// no unmarshalling
	return
}

// Update packets
type ServerDate struct { // Type 107
	CurrentDate uint32 // Current game date.
}

func (p ServerDate) String() string {
	return "SERVER_DATE"
}
func (p *ServerDate) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerClientJoin struct { // Type 108
	ID uint32 // ID of the new client.
}

func (p ServerClientJoin) String() string {
	return "SERVER_CLIENT_JOIN"
}
func (p *ServerClientJoin) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerClientInfo struct { // Type 109
	ID       uint32 // ID of the client.
	Address  string // Network address of the client.
	Name     string // Name of the client.
	Language uint8  // Language of the client.
	JoinDate uint32 // Date the client joined the game.
	Company  uint8  // ID of the company the client is playing as (255 for spectators).
}

func (p ServerClientInfo) String() string {
	return "SERVER_CLIENT_INFO"
}
func (p *ServerClientInfo) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerClientUpdate struct { // Type 110
	ID      uint32 // ID of the client.
	Name    string // Name of the client.
	Company uint8  // ID of the company the client is playing as (255 for spectators).
}

func (p ServerClientUpdate) String() string {
	return "SERVER_CLIENT_UPDATE"
}
func (p *ServerClientUpdate) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerClientQuit struct { // Type 111
	ID uint32 // ID of the leaving client.
}

func (p ServerClientQuit) String() string {
	return "SERVER_CLIENT_QUIT"
}
func (p *ServerClientQuit) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerClientError struct { // Type 112
	ID    uint32 // ID of the client throwing the error.
	Error uint8  // Error the client made (see NetworkErrorCode).
}

func (p ServerClientError) String() string {
	return "SERVER_CLIENT_ERROR"
}
func (p *ServerClientError) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerCompanyNew struct { // Type 113
	ID uint8 // ID of the new company.
}

func (p ServerCompanyNew) String() string {
	return "SERVER_COMPANY_NEW"
}
func (p *ServerCompanyNew) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
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

func (p ServerCompanyInfo) String() string {
	return "SERVER_COMPANY_INFO"
}
func (p *ServerCompanyInfo) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
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

func (p ServerCompanyUpdate) String() string {
	return "SERVER_COMPANY_UPDATE"
}
func (p *ServerCompanyUpdate) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerCompanyRemove struct { // Type 116
	ID     uint8 // ID of the company.
	Reason uint8 // Reason for being removed (see #AdminCompanyRemoveReason).
}

func (p ServerCompanyRemove) String() string {
	return "SERVER_COMPANY_REMOVE"
}
func (p *ServerCompanyRemove) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
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

func (p ServerCompanyEconomy) String() string {
	return "SERVER_COMPANY_ECONOMY"
}
func (p *ServerCompanyEconomy) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
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

func (p ServerCompanyStats) String() string {
	return "SERVER_COMPANY_STATS"
}
func (p *ServerCompanyStats) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerChat struct { // Type 119
	Action      uint8  // Action such as NETWORK_ACTION_CHAT_CLIENT (see #NetworkAction).
	Destination uint8  // Destination type such as DESTTYPE_BROADCAST (see #DestType).
	ID          uint32 // ID of the client who sent this message.
	Message     string // Message.
	Money       uint64 // Money (only when it is a 'give money' action).
}

func (p ServerChat) String() string {
	return "SERVER_CHAT"
}
func (p *ServerChat) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerRcon struct { // Type 120
	Colour uint8  // Colour as it would be used on the server or a client.
	Output string // Output of the executed command.
}

func (p ServerRcon) String() string {
	return "SERVER_RCON"
}
func (p *ServerRcon) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerConsole struct { // Type 121
	Origin  string // The origin of the text, e.g. "console" for console, or "net" for network related (debug) messages.
	Message string // Text as found on the console of the server.
}

func (p ServerConsole) String() string {
	return "SERVER_CONSOLE"
}
func (p *ServerConsole) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerCmdNames struct { // Type 122
	/*
	* NOTICE: Pack provided with this packet is not stable and will not be
	*         treated as such. Do not rely on IDs or names to be constant
	*         across different versions / revisions of OpenTTD.
	*         Pack provided in this packet is for logging purposes only.
	 */
	Commands map[uint16]string // Map of the ID of the DoCommand with the name of it.
}

func (p ServerCmdNames) String() string {
	return "SERVER_CMD_NAMES"
}

func (p *ServerCmdNames) Unmarshal(buffer *bytes.Buffer) (err error) {
	var next bool
	p.Commands = map[uint16]string{}
	binary.Read(buffer, binary.LittleEndian, &next)
	for next {
		// there are settings to read
		var cmdSlot uint16
		var cmdName string
		binary.Read(buffer, binary.LittleEndian, &cmdSlot)
		nvBytes, _ := buffer.ReadBytes(byte(0))
		cmdName = string(bytes.Trim(nvBytes, "\x00"))
		p.Commands[cmdSlot] = cmdName
		binary.Read(buffer, binary.LittleEndian, &next)
	}
	return nil
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

func (p ServerCmdLogging) String() string {
	return "SERVER_CMD_LOGGING"
}
func (p *ServerCmdLogging) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerGamescript struct { // Type 124
	// Pack on this isn't available in the source? Making an assumption.
	Json string // JSON string from the GameScript.
}

func (p ServerGamescript) String() string {
	return "SERVER_GAMESCRIPT"
}
func (p *ServerGamescript) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerRconEnd struct { // Type 125
	Command string // The command as requested by the admin connection.
}

func (p ServerRconEnd) String() string {
	return "SERVER_RCON_END"
}
func (p *ServerRconEnd) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}

type ServerPong struct { // Type 126
	Token uint32 // Integer value requested in the Ping.
}

func (p ServerPong) String() string {
	return "SERVER_PONG"
}
func (p *ServerPong) Unmarshal(buffer *bytes.Buffer) (err error) {
	return genericUnmarshal(p, buffer)
}
