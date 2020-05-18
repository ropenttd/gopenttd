package admin

import (
	"github.com/ropenttd/gopenttd/internal/helpers"
	"github.com/ropenttd/gopenttd/pkg/util"
	"net"
	"time"
)

type OpenttdExtendedCompany struct {
	// The company name.
	Name string `json:"name"`
	// The manager of the company.
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

type OpenttdClient struct {
	// The client nickname.
	Name string `json:"name"`
	// The client's IP address.
	Address net.IP `json:"address"`
	// The language the client is declaring themselves as.
	Language util.OpenttdLanguage `json:"language"`
	// The date the client joined the game.
	JoinDate time.Time `json:"date_join"`
	// The ID of the company that the client is playing in (set to 255 for spectators)
	Company uint8 `json:"company"`
}

type OpenttdExtendedServerState struct {
	// Status is set to True if the server was parsable, false otherwise.
	Status bool `json:"status"`
	// Error contains any parsing errors.
	Error error `json:"-"`

	// Dedicated is set if the server reports itself to be a Dedicated Server (instead of a Listen Server).
	Dedicated bool `json:"dedicated"`
	// Name is the server's advertised Hostname.
	Name string `json:"name"`
	// Version is the server's currently running Version string.
	Version string `json:"version"`
	// NeedPass defines whether the server is private (i.e needs a password)
	NeedPass bool `json:"need-pass"`

	// Language is the given Server Language.
	// If you need the string version of a language, use
	// gopenttd.GetISOLanguage(OpenttdServerState.Language) or gopenttd.GetLanguage(OpenttdServerState.Language)
	Language util.OpenttdLanguage `json:"language"`

	// Environment is the environment identifier.
	// If you need the string version of the environment, use
	// gopenttd.GetEnvironment(OpenttdServerState.Environment)
	Environment util.OpenttdEnvironment `json:"environment"`

	// Map is the name of the map currently running. Note that this will be set to "Random Map" if the map was generated by a seed.
	Map string `json:"map_name"`

	// Seed is the seed used to generate the map, if available (usually only through the admin port). nil if unavailable.
	Seed uint32 `json:"map_seed"`

	// MapHeight and MapWidth are the height and width of the current map in tiles.
	MapHeight uint16 `json:"map_height"`
	MapWidth  uint16 `json:"map_width"`

	// DateStart and DateCurrent are time objects relating to the start of the game and the current date.
	DateStart   time.Time `json:"date_start"`
	DateCurrent time.Time `json:"date_current"`

	// The following should be self explanatory.
	NumClients    int `json:"clients_active"`
	NumSpectators int `json:"spectators_active"`
	NumCompanies  int `json:"companies_active"`

	// Clients is a map of client IDs and client data.
	Clients map[uint32]OpenttdClient `json:"clients"`

	// Companies is a map of company IDs and company data.
	Companies map[uint8]OpenttdExtendedCompany `json:"companies"`
}

func (s *OpenttdExtendedServerState) updateCounts() {
	s.NumCompanies = len(s.Companies)
	s.NumClients = len(s.Clients)
	s.NumSpectators = 0
	for _, c := range s.Clients {
		if c.Company == 255 {
			lstate.NumSpectators += 1
		}
	}
}
