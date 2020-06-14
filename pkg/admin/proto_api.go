package admin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/admin/enum"
	"github.com/ropenttd/gopenttd/pkg/admin/packets"
	"github.com/ropenttd/gopenttd/pkg/util"
	"net"
	"time"
)

// Open creates a connection to the OpenTTD server.
func (s *Session) Open() error {
	s.log(LogInformational, "called")

	var err error

	// Prevent Open or other major Session functions from
	// being called while Open is still running.
	s.Lock()
	defer s.Unlock()

	// If the connection is already open, bail out here.
	if s.conn != nil {
		return ErrAlreadyConnected
	}

	// Connect to the server
	server := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
	serverAddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return err
	}

	s.log(LogInformational, "connecting to server %s", server)

	// Open the connection
	s.conn, err = net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		s.log(LogWarning, "error connecting to game %s, %s", server, err)
		s.conn = nil // Just to be safe.
		return err
	}

	defer func() {
		if err != nil {
			s.conn.Close()
			s.conn = nil
		}
	}()

	// We must first authenticate with the server before proceeding any further
	err = s.identify()
	if err != nil {
		err = fmt.Errorf("error sending identify packet to server: %s", err)
		return err
	}

	// Now OpenTTD should send us a Protocol message.
	mt, m, err := readPacketFromTcpConn(s.conn)
	if err != nil {
		return err
	}
	e, err := s.onEvent(mt, m)
	if err != nil {
		return err
	}
	if e.Type == fullEventType || e.Type == bannedEventType || e.Type == errorEventType {
		// The server doesn't want us for some reason
		s.log(LogError, "Server refused our connection: %s", e)
		return err
	}
	if e.Type != protocolEventType {
		// This is not fatal, but it does not follow the standard.
		s.log(LogWarning, "Expected PROTOCOL, instead got:\n%#v\n", e)
	}

	s.log(LogInformational, "We are now connected to OpenTTD, emitting connect event")
	s.handleEvent(connectEventType, &Connect{})

	// Create listening chan outside of listen, as it needs to happen inside the
	// mutex lock and needs to exist before calling heartbeat and listen
	// go routines.
	s.listening = make(chan interface{})

	// Start sending heartbeats and reading messages from the game.
	go s.heartbeat(s.conn, s.listening)
	go s.listen(s.conn, s.listening)

	s.log(LogInformational, "exiting")
	return nil
}

func (s *Session) Close() (err error) {
	s.Lock()
	s.Ready = false
	// Be polite, if we can
	if s.conn != nil {
		writePacketToTcpConn(s.conn, packets.AdminQuit{})
		// Close the connection
		s.conn.Close()
	}

	// Nil out the connection
	s.conn = nil

	s.log(LogInformational, "emit disconnect event")
	s.handleEvent(disconnectEventType, &Disconnect{})
	s.Unlock()

	return
}

// identify sends the authentication packet to the server
func (s *Session) identify() (err error) {

	data := packets.AdminJoin{
		Password:   s.Password,
		ClientName: s.UserAgent,
		Version:    VERSION,
	}

	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)

	return err
}

// isValidUpdateFrequency checks whether the given UpdateType can be requested at the given Frequency
// This requires valid data from the Protocol packet (i.e we have to be connected)
func (s *Session) isValidUpdateFrequency(t enum.UpdateType, f enum.UpdateFrequency) bool {
	if s.pollrates == nil {
		// We have no idea.
		return false
	}
	if v, ok := s.pollrates[t]; ok {
		// Bitwise check
		return v&uint16(f) != 0
	} else {
		return false
	}
}

// RequestUpdates sends a request to receive updates of the given type from the server at a given interval.
// Supplying a frequency of POLL is invalid and will return an error - use Session.Poll()
func (s *Session) RequestUpdates(t enum.UpdateType, f enum.UpdateFrequency) (err error) {
	// TODO cache these on the Session struct in case we need to reconnect to the server

	if f == enum.UpdateFrequencyPoll || !s.isValidUpdateFrequency(t, f) {
		// The server will ignore us or refuse us
		return ErrInvalidUpdateFrequency
	}

	data := packets.AdminUpdateFrequency{
		Type:      t,
		Frequency: f,
	}
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)
	return err
}

// Poll sends a request to receive one update for the given UpdateType and ID.
// If you can't poll for the given UpdateType, this returns an error.
func (s *Session) Poll(t enum.UpdateType, id uint32) (err error) {
	if !s.isValidUpdateFrequency(t, enum.UpdateFrequencyPoll) {
		// We can't poll for this thing
		return ErrInvalidUpdateFrequency
	}
	data := packets.AdminPoll{
		Type: t,
		ID:   ^uint32(0),
	}
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)
	return err
}

// Chat sends a chat message (who'dve thought it?)
func (s *Session) Chat(act enum.Action, dest enum.Destination, destID uint32, message string) (err error) {
	data := packets.AdminChat{
		Action:        act,
		Destination:   dest,
		DestinationID: destID,
		Message:       message,
	}
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)
	return err
}

// Rcon sends a non-blocking RCON command to the server.
// You are expected to watch for events of type Rcon and RconEnd to determine the result if you use this.
func (s *Session) Rcon(com string) (err error) {
	data := packets.AdminRcon{
		Command: com,
	}
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)
	return err
}

// GamescriptCommand sends a non-blocking Gamescript command to the server.
// You are expected to watch for events of type Gamescript to determine the result if you use this.
func (s *Session) GamescriptCommand(json string) (err error) {
	data := packets.AdminGamescript{
		Json: json,
	}
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)
	return err
}

// listen polls the websocket connection for events, it will stop when the
// listening channel is closed, or an error occurs.
func (s *Session) listen(conn *net.TCPConn, listening <-chan interface{}) {

	s.log(LogInformational, "called")

	for {

		messageType, message, err := readPacketFromTcpConn(conn)

		if err != nil {

			// Detect if we have been closed manually. If a Close() has already
			// happened, the socket we are listening on will be different to
			// the current session.
			s.RLock()
			sameConnection := s.conn == conn
			s.RUnlock()

			if sameConnection {

				s.log(LogWarning, "error reading from game %s, %s", s.Hostname, err)
				// There has been an error reading, close the socket so that
				// OnDisconnect event is emitted.
				err := s.Close()
				if err != nil {
					s.log(LogWarning, "error closing session connection, %s", err)
				}

				s.log(LogInformational, "calling reconnect() now")
				s.reconnect()
			}

			return
		}

		select {

		case <-listening:
			return

		default:
			s.onEvent(messageType, message)

		}
	}
}

func (s *Session) reconnect() {

	s.log(LogInformational, "called")

	var err error

	if s.ShouldReconnectOnError {

		wait := time.Duration(1)

		for {
			s.log(LogInformational, "trying to reconnect to game")

			err = s.Open()
			if err == nil {
				s.log(LogInformational, "successfully reconnected to game")
			}

			// Certain race conditions can call reconnect() twice. If this happens, we
			// just break out of the reconnect loop
			if err == ErrAlreadyConnected {
				s.log(LogInformational, "Connection already active, no need to reconnect")
				return
			}

			s.log(LogError, "error reconnecting to game, %s", err)

			<-time.After(wait * time.Second)
			wait *= 2
			if wait > 600 {
				wait = 600
			}
		}
	}
}

// FailedPongs is the Number of pong intervals to wait until forcing a connection restart.
const FailedPongs = 5 * time.Millisecond

// HeartbeatLatency returns the latency between heartbeat acknowledgement and heartbeat send.
func (s *Session) HeartbeatLatency() time.Duration {

	return s.LastPong.Sub(s.LastPing)

}

// heartbeat sends regular heartbeats to OpenTTD to ensure the server is still available.
func (s *Session) heartbeat(conn *net.TCPConn, listening <-chan interface{}) {

	s.log(LogInformational, "called")

	if listening == nil || conn == nil {
		return
	}

	heartbeatIntervalMsec := time.Duration(10000)

	var err error
	ticker := time.NewTicker(heartbeatIntervalMsec * time.Second)
	defer ticker.Stop()

	for {
		s.RLock()
		last := s.LastPong
		s.RUnlock()
		s.log(LogDebug, "sending game ping")
		s.connMutex.Lock()
		s.LastPing = time.Now().UTC()
		// very lazy implementation of token for very lazy people
		err = writePacketToTcpConn(conn, packets.AdminPing{Token: 1})
		s.connMutex.Unlock()
		if err != nil || time.Now().UTC().Sub(last) > (heartbeatIntervalMsec*FailedPongs) {
			if err != nil {
				s.log(LogError, "error sending heartbeat to server %s, %s", s.Hostname, err)
			} else {
				s.log(LogError, "haven't gotten a pong in %v, triggering a reconnection", time.Now().UTC().Sub(last))
			}
			s.Close()
			s.reconnect()
			return
		}
		s.Lock()
		s.Ready = true
		s.Unlock()

		select {
		case <-ticker.C:
			// continue loop and send heartbeat
		case <-listening:
			return
		}
	}
}

// onEvent is the "event handler" for all messages received on the
// OpenTTD Admin connection.
//
// If you use the AddHandler() function to register a handler for a
// specific event this function will pass the event along to that handler.
//
// If you use the AddHandler() function to register a handler for the
// "OnEvent" event then all events will be passed to that handler.
func (s *Session) onEvent(messageType uint8, message []byte) (*Event, error) {

	var err error

	// Pack the event into an Event struct.
	var e = new(Event)
	e.Type = messageType
	e.RawData = message

	s.log(LogDebug, "Type: %d, Data: %s\n\n", e.Type, string(e.RawData))

	if e.Type == packetIndexServerPong {
		s.Lock()
		s.LastPong = time.Now().UTC()
		s.Unlock()
		s.log(LogDebug, "got pong")
		return e, nil
	}

	// Map event to registered event handlers and pass it along to any registered handlers.
	if eh, ok := registeredInterfaceProviders[e.Type]; ok {
		e.Struct = eh.New()

		// Attempt to unmarshal our event.
		if err = ottdUnmarshal(e.RawData, e.Struct); err != nil {
			s.log(LogError, "error unmarshalling %s event, %s", e.Type, err)
		}

		// Send event to any registered event handlers for its type.
		s.handleEvent(e.Type, e.Struct)
	} else {
		s.log(LogWarning, "unknown event: Type: %d, Data: %s", e.Type, string(e.RawData))
	}

	return e, nil
}

// readPacket is a non-public packet reader.
func readPacketFromTcpConn(r *net.TCPConn) (messageType uint8, p []byte, err error) {
	// Read the first part
	lengthBytes := make([]byte, 2)
	_, err = r.Read(lengthBytes)
	if err != nil {
		return messageType, p, err
	}
	if len(lengthBytes) < 2 {
		return messageType, p, errors.New("received a packet shorter than the required length")
	}

	packLength := int(binary.LittleEndian.Uint16(lengthBytes))

	data := make([]byte, packLength-2)

	readLen, err := r.Read(data)
	if err != nil {
		return messageType, p, err
	}

	if readLen+2 != packLength {
		// ignore the packet
		return messageType, p, errors.New(fmt.Sprint("invalid reported buffer length: got ", readLen, ", expected ", packLength))
	}

	messageType = data[0]

	// return the bytes AFTER the messageType as a byte array for compatibility
	return messageType, data[1:], err
}

func writePacketToTcpConn(c *net.TCPConn, packet packets.AdminRequestPacket) (err error) {
	// A packet has fallen into the writer in gopenttd city!
	// Start the new packet builder!

	// HEY

	// Build the packet
	data := packet.Pack()
	msg := new(bytes.Buffer)

	// And off to the rescue
	// Write the packet length header
	// Length is +3 because of the metadata fields we're adding at the beginning
	msgLength := uint16(data.Len() + 3)
	binary.Write(msg, binary.LittleEndian, msgLength)

	// Lower the data line
	msg.WriteByte(uint8(packet.PacketType()))
	msg.Write(data.Bytes())

	// And make the write
	sendLen, err := c.Write(msg.Bytes())
	if err != nil {
		return err
	}
	if sendLen != int(msgLength) {
		return util.ErrBadWrite
	}

	// The new packet building collection from gopenttd.
	return nil
}
