// Package gopenttd_admin provides primitives for querying OpenTTD game server admin interfaces.
package gopenttd

import (
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/consts"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

func (c *OpenttdAdminConnection) Open() (err error) {
	if c.Connection != nil || c.connected {
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
	return c.Auth()
}

func (c *OpenttdAdminConnection) Auth() (err error) {
	log.Info("Authenticating with ", c.Hostname)
	pack := packets.AdminJoin{
		Password:   c.Password,
		ClientName: c.ClientName,
		Version:    "testing",
	}
	c.Writer.Write(pack)

	// Watch the first packet ourselves (yes I'm aware this is a little hacky)
	packet, err := c.Reader.Read()

	if err == io.EOF {
		// Some kind of connection error
		return errAuthentication
	} else if err != nil {
		log.Error("Read error during connection: ", err)
		return err
	}

	log.Info("Authenticated successfully")
	c.connected = true

	c.publish(packet)
	return
}

func (c *OpenttdAdminConnection) Close() (err error) {
	c.Writer.Write(packets.AdminQuit{})
	log.Infof("Connection to %s:%d closed", c.Hostname, c.Port)
	c.connected = false
	return c.Connection.Close()
}

func (c *OpenttdAdminConnection) Watch() (err error) {
	if !c.connected {
		return errNotConnected
	}
	for {
		// Start a loop to read from TCP, and exit out if an error or timeout is experienced
		packet, err := c.Reader.Read()

		if err != nil {
			// Some kind of connection error
			log.Error("Read error: ", err)
			return err
		}

		c.publish(packet)
	}
	return nil
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
func ScanServerAdm(host string, port int, password string) (err error) {
	obj := OpenttdAdminConnection{Hostname: host, Port: port, Password: password, ClientName: "gopenttd"}

	ch := make(chan packets.AdminResponsePacket)
	obj.Subscribe(ch)

	err = obj.Open()
	if err != nil {
		return err
	}
	defer obj.Close()
	go obj.Watch()

	obj.Writer.Write(packets.AdminPoll{
		Type: consts.UpdateTypeDate,
		ID:   ^uint32(0),
	})

	for {
		select {
		case d := <-ch:
			log.Info(d)
		}
	}

	return
}
