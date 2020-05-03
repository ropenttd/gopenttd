package openttd_consts_admin

// As defined in https://github.com/OpenTTD/OpenTTD/blob/master/src/network/core/tcp_admin.h

const UpdateTypeDate = 0
const UpdateTypeClientInfo = 1
const UpdateTypeCompanyInfo = 2
const UpdateTypeCompanyEconomy = 3
const UpdateTypeCompanyStats = 4
const UpdateTypeChat = 5
const UpdateTypeConsole = 6
const UpdateTypeCmdNames = 7
const UpdateTypeCmdLogging = 8

const UpdateFrequencyPoll = 0x01
const UpdateFrequencyDaily = 0x02
const UpdateFrequencyWeekly = 0x04
const UpdateFrequencyMonthly = 0x08
const UpdateFrequencyQuarterly = 0x10
const UpdateFrequencyAnnually = 0x20
const UpdateFrequencyAutomatically = 0x40

const CompanyRemoveReasonManual = 0
const CompanyRemoveReasonAutoclean = 1
const CompanyRemoveReasonBankrupt = 2

const ActionJoin = 0x00
const ActionLeave = 0x01
const ActionServerMessage = 0x02
const ActionChat = 0x03
const ActionChatCompany = 0x04
const ActionChatClient = 0x05
const ActionGiveMoney = 0x06
const ActionNameChange = 0x07
const ActionCompanySpectator = 0x08
const ActionCompanyJoin = 0x09
const ActionCompanyNew = 0x0A

const DestinationBroadcast = 0x00
const DestinationTeam = 0x01
const DestinationClient = 0x02

const NetErrorGeneral = 0x00

// Signals from clients
const NetErrorDesync = 0x01
const NetErrorSavegameFailed = 0x02
const NetErrorConnectionLost = 0x03
const NetErrorIllegalPacket = 0x04
const NetErrorNewgrfMismatch = 0x05

// Signals from servers
const NetErrorNotAuthorized = 0x06
const NetErrorNotExpected = 0x07
const NetErrorWrongRevision = 0x08
const NetErrorNameInUse = 0x09
const NetErrorWrongPassword = 0x0A
const NetErrorCompanyMismatch = 0x0B
const NetErrorKicked = 0x0C
const NetErrorCheater = 0x0D
const NetErrorFull = 0x0E
const NetErrorTooManyCommands = 0x0F
