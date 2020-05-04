// Package gopenttd_admin provides primitives for querying OpenTTD game server admin interfaces.
package gopenttd_admin

import (
	"fmt"
	"github.com/ropenttd/gopenttd/internal/openttd_consts_admin"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

func (c *OpenttdAdminConnection) Open() (err error) {
	if c.Connection != nil {
		// We're already connected
		return nil
	}
	// Determine the correct TCP address
	server := fmt.Sprintf("%s:%d", c.Hostname, c.Port) // "255.255.255.255", 10000
	serverAddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return err
	}

	// Open the connection
	c.Connection, err = net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return err
	}
	c.Reader = NewPacketReader(c.Connection)
	c.Writer = NewPacketWriter(c.Connection)

	log.Info("Admin connection open to ", server)
	c.SendAuth()
	return nil
}

func (c *OpenttdAdminConnection) SendAuth() (err error) {
	log.Info("Authenticating with ", c.Hostname)
	pack := openttd_consts_admin.AdminJoin{
		Password:   c.Password,
		ClientName: c.ClientName,
		Version:    "0.0.1",
	}
	c.Writer.Write(pack)
	return
}

func (c *OpenttdAdminConnection) Close() (err error) {
	c.Writer.Write(openttd_consts_admin.AdminQuit{})
	log.Debugf("Connection to %s:%d closed", c.Hostname, c.Port)
	return c.Connection.Close()
}

func (c *OpenttdAdminConnection) Watch() (err error) {
	for {
		// Start a loop to read from TCP, and exit out if an error or timeout is experienced
		packet, err := c.Reader.Read()

		if err != nil && err != io.EOF {
			// Some kind of connection error
			log.Error("Read error: ", err)
			return err
		}

		switch packet.(type) {
		case openttd_consts_admin.ServerProtocol:
			log.Debugf("Received ServerProtocol packet %v", packet)
		case openttd_consts_admin.ServerWelcome:
			log.Debugf("Received ServerWelcome packet %v", packet)
		default:
			log.Infof("Received Unknown packet: %v", packet)
		}
	}
	return nil
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
// Connections time out after 10 seconds
func ScanServerAdm(host string, port int, password string) (err error) {
	obj := OpenttdAdminConnection{Hostname: host, Port: port, Password: password, ClientName: "gopenttd"}
	err = obj.Open()
	if err != nil {
		return err
	}
	defer obj.Close()
	obj.Watch()
	return
}
