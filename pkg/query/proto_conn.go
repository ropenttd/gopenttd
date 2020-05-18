package query

import "net"

type GameConn interface {
	Open() error
	Close() error
	Query() error
}

type OpenttdClientConnection struct {
	Hostname   string
	Port       int
	Connection *net.UDPConn
}
