package admin

import (
	"github.com/ropenttd/gopenttd/internal/helpers"
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
	"github.com/ropenttd/gopenttd/pkg/util"
	"net"
	"sync"
	"time"
)

// A Session represents a connection to the OpenTTD Admin protocol.
type Session struct {
	sync.RWMutex

	// General configurable settings.

	// Hostname to connect to.
	Hostname string

	// Port to connect to (note this is the admin port, usually 3977: not the game port!)
	Port int

	// Authentication token for this session
	Password string

	// Debug for printing JSON request/responses
	LogLevel int

	// Should the session reconnect on errors.
	ShouldReconnectOnError bool

	// Should state tracking be enabled.
	// State tracking automatically updates the State object
	// when events occur - if you're purely writing real time stuff,
	// you can turn this off.
	StateEnabled bool

	// Whether or not to call event handlers synchronously.
	// e.g false = launch event handlers in their own goroutines.
	SyncEvents bool

	// Exposed but should not be modified by User.

	// Whether the connection is ready
	Ready bool

	// Managed state object, updated internally with events when
	// StateEnabled is true.
	State *State

	// The user agent
	UserAgent string

	// Stores the last Pong that was recieved (in UTC)
	LastPong time.Time

	// Stores the last Heartbeat sent (in UTC)
	LastPing time.Time

	// Event handlers
	handlersMu   sync.RWMutex
	handlers     map[uint8][]*eventHandlerInstance
	onceHandlers map[uint8][]*eventHandlerInstance

	// The TCP connection.
	conn *net.TCPConn

	// Acceptable polling rates
	pollrates map[enum.UpdateType]uint16

	// When nil, the session is not listening.
	listening chan interface{}

	// used to make sure writes do not happen concurrently
	connMutex sync.Mutex
}

type Company struct {
	// The company name.
	Name string `json:"name"`
	// The manager of the company (client-defined, not definitively the name of a client - it appears under the portrait).
	Manager string `json:"manager"`
	// The colour representing this company in-game.
	Colour helpers.OpenttdColour `json:"colour"`
	// The year the company was first founded.
	YearStart uint32 `json:"start_year"`
	// The value of the company, in credits.
	Value uint64 `json:"value"`
	// The amount of disposable cash the company has.
	Money uint64 `json:"cash"`
	// The company's current income. Can go negative if they're making a loss.
	Income int64 `json:"income"`
	// The company's current loan.
	Loan uint64 `json:"loan"`
	// The number of quarters this company has been in a bankruptcy state (this may mean they're about to be dissolved!)
	Bankruptcy uint8 `json:"bankruptcy"`

	// Share data
	// Share 1 is owned by...
	Share1 uint8 `json:"share1"`
	// Share 2 is owned by...
	Share2 uint8 `json:"share2"`
	// Share 3 is owned by...
	Share3 uint8 `json:"share3"`
	// Share 4 is owned by...
	Share4 uint8 `json:"share4"`

	// Cargo delivered in the current quarter.
	CargoThisQuarter uint16 `json:"cargo_current"`
	// Cargo delivered in the last full quarter.
	CargoLastQuarter uint16 `json:"cargo_last"`
	// Cargo delivered in the quarter before the last full one.
	CargoPreviousQuarter uint16 `json:"cargo_previous"`
	// Company value in the last quarter.
	ValueLastQuarter uint64 `json:"value_last"`
	// Company value in the previous quarter.
	ValuePreviousQuarter uint64 `json:"value_previous"`
	// Company performance index in the last quarter, out of 1000.
	PerformanceLastQuarter uint16 `json:"performance_last"`
	// Company performance index in the previous quarter.
	PerformancePreviousQuarter uint16 `json:"performance_previous"`

	// Whether the company has a password set.
	Passworded bool `json:"is_passworded"`
	// Whether the company is controlled by an AI.
	AI bool `json:"is_ai"`

	// A count of the vehicles and stations the company has.
	Vehicles util.OpenttdTypeCounts `json:"vehicle_count"`
	Stations util.OpenttdTypeCounts `json:"station_count"`
}

type Client struct {
	// The client nickname.
	Name string `json:"name"`
	// The client's IP address.
	// This can be nil if the client is the host or is connecting via an invalid address, so watch out!
	Address net.IP `json:"address"`
	// The language the client is declaring themselves as.
	Language util.OpenttdLanguage `json:"language"`
	// The date the client joined the game.
	// This is set to zero if the client is the server host, so watch out!
	JoinDate time.Time `json:"date_join"`
	// The ID of the company that the client is playing in (set to 255 for spectators)
	Company uint8 `json:"company"`
}
