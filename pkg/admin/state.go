package admin

import (
	"github.com/ropenttd/gopenttd/internal/helpers"
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
	"github.com/ropenttd/gopenttd/pkg/util"
	"net"
	"sync"
	"time"
)

// A State contains the current known state of the server.
type State struct {
	sync.RWMutex

	// Dedicated is set if the server reports itself to be a Dedicated Server (instead of a Listen Server).
	Dedicated bool `json:"dedicated"`
	// Name is the server's advertised Hostname.
	Name string `json:"name"`
	// Version is the server's currently running Version string.
	Version string `json:"version"`
	// ProtocolVersion is the server's protocol version (i.e what feature set it supports).
	ProtocolVersion uint8 `json:"protocol_version"`
	// NeedPass defines whether the server is private (i.e needs a password)
	NeedPass bool `json:"need-pass"`

	// Language is the given Server Language.
	// If you need the string version of a language, use
	// gopenttd.GetISOLanguage(OpenttdServerState.Language) or gopenttd.GetLanguage(OpenttdServerState.Language)
	Language util.OpenttdLanguage `json:"language"`

	// Landscape is the environment identifier.
	// If you need the string version of the environment, use
	// gopenttd.GetEnvironment(OpenttdServerState.Environment)
	Landscape util.OpenttdEnvironment `json:"landscape"`

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

	// If you want a count of the number of companies and clients connected,
	// you can either count Clients and Companies yourself,
	// or call State.Counts()

	// Clients is a map of client IDs and client data.
	Clients map[uint32]Client `json:"clients"`

	// Companies is a map of company IDs and company data.
	Companies map[uint8]Company `json:"companies"`
}

// NewState creates an empty state.
func NewState() *State {
	return &State{
		Clients:   map[uint32]Client{},
		Companies: map[uint8]Company{},
	}
}

// Counts simply counts the number of companies, spectators, and clients connected.
// Note that Clients is a count of ALL clients, INCLUDING spectators. (i.e Clients will always be bigger than Spectators)
func (s *State) Counts() (clients int, spectators int, companies int) {
	for _, c := range s.Clients {
		clients++
		if c.Company == 255 {
			spectators++
		}
	}
	companies = len(s.Companies)
	return
}

// OnProtocol takes a Protocol event and updates state with the protocol version and settings.
func (s *State) onProtocol(se *Session, r *Protocol) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	s.ProtocolVersion = r.Version

	// Reset polling rates
	se.pollrates = map[enum.UpdateType]uint16{}
	for k, v := range r.Settings {
		se.pollrates[enum.UpdateType(k)] = v
	}
	return
}

// OnWelcome takes a Welcome event and updates all internal state.
func (s *State) onWelcome(se *Session, r *Welcome) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	s.Name = r.Name
	s.Version = r.Version
	s.Dedicated = r.Dedicated
	s.Map = r.Map
	s.Seed = r.Seed
	s.Landscape = util.OpenttdEnvironment(r.Landscape)
	s.DateStart = util.DateFormat(r.StartDate)
	s.MapWidth = r.MapWidth
	s.MapHeight = r.MapHeight

	return nil
}

// OnDate takes a Date event and updates the current date in state.
func (s *State) onDate(se *Session, r *Date) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	s.DateCurrent = util.DateFormat(r.CurrentDate)
	return
}

// OnClientJoin just requests further information about the client.
func (s *State) onClientJoin(se *Session, r *ClientJoin) (err error) {
	if s == nil || s.Clients == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.Clients[r.ID]; ok {
		// client in clients state, wtf?
		se.log(LogWarning, "Joining Client %d appears to be in the state already - this shouldn't happen?", r.ID)
	}

	// If they already exist, we give them a new state anyway
	s.Clients[r.ID] = Client{}

	// Poll for more information about this client
	err = se.Poll(enum.UpdateTypeClientInfo, r.ID)
	return err
}

// OnClientInfo updates the information relating to the client that is referenced.
func (s *State) onClientInfo(se *Session, r *ClientInfo) (err error) {
	if s == nil || s.Clients == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	var cli Client
	if res, ok := s.Clients[r.ID]; ok {
		// client in clients state
		cli = res
	} else {
		// client not in clients state
		cli = Client{}
	}

	cli.Name = r.Name
	cli.Language = util.OpenttdLanguage(r.Language)
	cli.JoinDate = util.DateFormat(r.JoinDate)
	cli.Address = net.ParseIP(r.Address)
	cli.Company = r.Company

	s.Clients[r.ID] = cli

	return err
}

// OnClientUpdate fires when the client updates either their name or their company.
func (s *State) onClientUpdate(se *Session, r *ClientUpdate) (err error) {
	if s == nil || s.Clients == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	var cli Client
	if res, ok := s.Clients[r.ID]; ok {
		// client in clients state
		cli = res
	} else {
		// client not in clients state
		cli = Client{}
	}

	cli.Name = r.Name
	cli.Company = r.Company

	s.Clients[r.ID] = cli

	return err
}

// OnClientQuit removes the client from the state.
func (s *State) onClientQuit(se *Session, r *ClientQuit) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.Clients[r.ID]; ok {
		delete(s.Clients, r.ID)
	} else {
		se.log(LogWarning, "Leaving Client %d does not appear to be in the state, ignoring", r.ID)
	}

	return err
}

// OnCompanyNew just requests further information about the company.
func (s *State) onCompanyNew(se *Session, r *CompanyNew) (err error) {
	if s == nil || s.Companies == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.Companies[r.ID]; ok {
		// company in state, we obviously missed the removal
		se.log(LogInformational, "New Company %d appears to be in the state already - we obviously missed the memo", r.ID)
	}

	// If they already exist, this gives them a new state anyway (which will be populated by the poll)
	s.Companies[r.ID] = Company{}

	// Poll for more information about this company
	// we have to cast our uint8 to a uint32 because that's what the poll packet expects
	err = se.Poll(enum.UpdateTypeCompanyInfo, uint32(r.ID))
	return err
}

// OnCompanyInfo updates the given company with the given information.
func (s *State) onCompanyInfo(se *Session, r *CompanyInfo) (err error) {
	if s == nil || s.Companies == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	var com Company
	if res, ok := s.Companies[r.ID]; ok {
		// company in state
		com = res
	} else {
		// company not in state
		com = Company{}
	}

	com.Name = r.Name
	com.Manager = r.Manager
	com.Colour = helpers.OpenttdColour(r.Colour)
	com.Passworded = r.Password
	// YearStart is an INTEGER containing the year of founding, NOT a time.Time
	com.YearStart = r.StartDate
	com.AI = r.IsAI

	s.Companies[r.ID] = com

	return err
}

// OnCompanyUpdate when the company updates something about their state.
func (s *State) onCompanyUpdate(se *Session, r *CompanyUpdate) (err error) {
	if s == nil || s.Companies == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	var com Company
	if res, ok := s.Companies[r.ID]; ok {
		// company in state
		com = res
	} else {
		// company not in state
		com = Company{}
	}

	com.Name = r.Name
	com.Manager = r.Manager
	com.Colour = helpers.OpenttdColour(r.Colour)
	com.Passworded = r.Password
	com.Bankruptcy = r.BankruptcyQuarters

	if se.State.ProtocolVersion <= 2 {
		// Company shares were removed in OpenTTD 14.0.
		com.Share1 = r.Share1
		com.Share2 = r.Share2
		com.Share3 = r.Share3
		com.Share4 = r.Share4
	}

	s.Companies[r.ID] = com

	return err
}

// OnCompanyRemove removes the company from the state.
func (s *State) onCompanyRemove(se *Session, r *CompanyRemove) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.Companies[r.ID]; ok {
		delete(s.Companies, r.ID)
	} else {
		se.log(LogWarning, "Dissolved Company %d does not appear to be in the state, ignoring", r.ID)
	}

	return err
}

// OnCompanyEconomy updates the given company with economy statistics.
func (s *State) onCompanyEconomy(se *Session, r *CompanyEconomy) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	var com Company
	if res, ok := s.Companies[r.ID]; ok {
		// company in state
		com = res
	} else {
		// company not in state
		com = Company{}
	}

	com.Money = r.Money
	com.Loan = r.Loan
	com.Income = r.Income
	com.CargoThisQuarter = r.CargoThisQuarter
	com.ValueLastQuarter = r.ValueLastQuarter
	com.PerformanceLastQuarter = r.PerformanceLastQuarter
	com.CargoLastQuarter = r.CargoLastQuarter
	com.ValuePreviousQuarter = r.ValuePreviousQuarter
	com.PerformancePreviousQuarter = r.PerformancePreviousQuarter
	com.CargoPreviousQuarter = r.CargoPreviousQuarter

	s.Companies[r.ID] = com

	return err

}

// OnCompanyStats updates the given company's vehicle counts.
func (s *State) onCompanyStats(se *Session, r *CompanyStats) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	var com Company
	if res, ok := s.Companies[r.ID]; ok {
		// company in state
		com = res
	} else {
		// company not in state
		com = Company{Vehicles: util.OpenttdTypeCounts{}, Stations: util.OpenttdTypeCounts{}}
	}

	com.Vehicles.Train = r.Trains
	com.Vehicles.Truck = r.Lorries
	com.Vehicles.Bus = r.Buses
	com.Vehicles.Aircraft = r.Planes
	com.Vehicles.Ship = r.Ships

	com.Stations.Train = r.TrainStations
	com.Stations.Truck = r.LorryStations
	com.Stations.Bus = r.BusStops
	com.Stations.Aircraft = r.Airports
	com.Stations.Ship = r.Harbours

	s.Companies[r.ID] = com

	return err

}

// OnInterface handles all events related to states.
func (s *State) OnInterface(se *Session, i interface{}) (err error) {
	if s == nil {
		return ErrNilState
	}

	switch r := i.(type) {
	case *Protocol:
		return s.onProtocol(se, r)
	case *Welcome:
		return s.onWelcome(se, r)
	}

	if !se.StateEnabled {
		return err
	}

	// State changes
	switch r := i.(type) {
	case *Date:
		err = s.onDate(se, r)
	case *ClientJoin:
		err = s.onClientJoin(se, r)
	case *ClientInfo:
		err = s.onClientInfo(se, r)
	case *ClientUpdate:
		err = s.onClientUpdate(se, r)
	case *ClientQuit:
		err = s.onClientQuit(se, r)
	case *CompanyNew:
		err = s.onCompanyNew(se, r)
	case *CompanyInfo:
		err = s.onCompanyInfo(se, r)
	case *CompanyUpdate:
		err = s.onCompanyUpdate(se, r)
	case *CompanyRemove:
		err = s.onCompanyRemove(se, r)
	case *CompanyEconomy:
		err = s.onCompanyEconomy(se, r)
	case *CompanyStats:
		err = s.onCompanyStats(se, r)
	}
	return err
}
