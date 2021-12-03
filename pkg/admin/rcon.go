package admin

import "github.com/ropenttd/gopenttd/pkg/admin/packets"

// RCON related stuff is dealt with in this file to help keep things a little tidier.

type rconRequest struct {
	Command      string `json:"command"`
	responseChan chan []Rcon
}

type rconResp struct {
	rcon    *Rcon
	rconEnd *RconEnd
}

// Rcon sends a non-blocking RCON command to the server.
// Use this when you don't care what the result is - if you do, use RconSync(command).
func (s *Session) Rcon(command string) (err error) {
	// we have to add this to the queue because the handleRconRequests queue will get out of step with commands otherwise

	obj := rconRequest{Command: command}
	s.rconQueue <- &obj
	return nil
}

// RconSync sends a blocking RCON command to the server, waits for a response, then returns a set of response packets.
// Please note: This will block your thread until we get a complete response from the server!
// If you don't care about the result, use Rcon(command).
func (s *Session) RconSync(command string) (ret []Rcon, err error) {
	rchan := make(chan []Rcon)
	obj := rconRequest{Command: command, responseChan: rchan}
	s.rconQueue <- &obj
	// Block on a response
	ret = <-obj.responseChan
	return ret, nil
}

func (s *Session) sendRconCommand(command string) (err error) {
	data := packets.AdminRcon{
		Command: command,
	}
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	err = writePacketToTcpConn(s.conn, data)
	return err
}

func (s *Session) handleRconRequests(listening <-chan interface{}) {

	s.log(LogDebug, "called")

	for {
		select {
		case <-listening:
			return
		case cmd := <-s.rconQueue:
			// Send it
			err := s.sendRconCommand(cmd.Command)
			if err != nil {
				if cmd.responseChan != nil {
					cmd.responseChan <- []Rcon{}
				}
				continue
			}
			var data []Rcon
			var run = true
			for run {
				v := <-s.rconChan
				switch {
				case v.rcon != nil:
					// not an ending packet
					data = append(data, *v.rcon)
				case v.rconEnd != nil:
					// ending packet
					if v.rconEnd.Command == cmd.Command {
						run = false
					}
				}
			}
			if cmd.responseChan != nil {
				cmd.responseChan <- data
			}
		}
	}
}
