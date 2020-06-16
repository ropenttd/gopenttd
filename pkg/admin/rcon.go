package admin

type rconRequest struct {
	Command      string `json:"command"`
	responseChan chan []Rcon
}

type rconResp struct {
	rcon    *Rcon
	rconEnd *RconEnd
}

func (s *Session) RconSync(command string) (ret []Rcon, err error) {
	// is this too wide of a scope to lock? we could have an RconLock otherwise
	s.Lock()
	rchan := make(chan []Rcon)
	obj := rconRequest{Command: command, responseChan: rchan}
	s.rconQueue = append(s.rconQueue, obj)
	s.Unlock()
	// Block on a response
	ret = <-obj.responseChan
	return ret, nil
}

func (s *Session) handleRconRequests(listening <-chan interface{}) {

	s.log(LogInformational, "called")

	for {
		var cmd rconRequest
		if len(s.rconQueue) == 0 {
			// no requests available right now
			break
		}

		// is this too wide of a scope to lock? we could have an RconLock otherwise
		s.Lock()

		// Pop the last thing on the stack
		cmd, s.rconQueue = s.rconQueue[len(s.rconQueue)-1], s.rconQueue[:len(s.rconQueue)-1]

		// Send it
		err := s.Rcon(cmd.Command)
		if err != nil {
			cmd.responseChan <- []Rcon{}
			break
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
		select {
		case <-listening:
			s.Unlock()
			return
		default:
			cmd.responseChan <- data
		}
		s.Unlock()
	}
}
