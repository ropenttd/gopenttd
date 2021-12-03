package enum

// As defined in https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_admin.h

type UpdateType uint8

const (
	UpdateTypeDate UpdateType = iota
	UpdateTypeClientInfo
	UpdateTypeCompanyInfo
	UpdateTypeCompanyEconomy
	UpdateTypeCompanyStats
	UpdateTypeChat
	UpdateTypeConsole
	UpdateTypeCmdNames
	UpdateTypeCmdLogging
	UpdateTypeGamescript
)

type UpdateFrequency uint16

const (
	UpdateFrequencyPoll    UpdateFrequency = 0x01 << iota
	UpdateFrequencyDaily                   // 0x02
	UpdateFrequencyWeekly                  // 0x04
	UpdateFrequencyMonthly                 // 0x08, etc
	UpdateFrequencyQuarterly
	UpdateFrequencyAnnually
	UpdateFrequencyAutomatically
)

type CompanyRemoveReason uint8

const (
	CompanyRemoveReasonManual CompanyRemoveReason = iota
	CompanyRemoveReasonAutoclean
	CompanyRemoveReasonBankrupt
)

type Action uint8

const (
	ActionJoin Action = 0x00
	ActionLeave
	ActionServerMessage
	ActionChat
	ActionChatCompany
	ActionChatClient
	ActionGiveMoney
	ActionNameChange
	ActionCompanySpectator
	ActionCompanyJoin
	ActionCompanyNew // 0x0A
)

type Destination uint8

const (
	DestinationBroadcast Destination = 0x00 << iota // All destinations
	DestinationTeam                                 // A specific team
	DestinationClient                               // A specific client
)

type NetError uint8

const (
	NetErrorGeneral NetError = 0x00 << iota // A general network failure

	// Signals from clients
	NetErrorDesync
	NetErrorSavegameFailed
	NetErrorConnectionLost
	NetErrorIllegalPacket
	NetErrorNewgrfMismatch // 0x05

	// Signals from servers
	NetErrorNotAuthorized
	NetErrorNotExpected
	NetErrorWrongRevision
	NetErrorNameInUse
	NetErrorWrongPassword
	NetErrorCompanyMismatch
	NetErrorKicked
	NetErrorCheater
	NetErrorFull
	NetErrorTooManyCommands // 0x0F
)

type ClientID uint32

const (
	// Client is not part of anything
	ClientIDInvalid = 0x00 << iota
	// Server is guaranteed to have this Client ID
	ClientIDServer = 0x01
	// The first Client ID
	ClientIDFirst = 0x02
)
