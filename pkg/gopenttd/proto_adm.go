// Package gopenttd provides primitives for querying OpenTTD game servers.
package gopenttd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ropenttd/gopenttd/internal/openttd_consts_admin"
	"net"
	"time"

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
	log.Info("Admin connection open to ", server)
	c.Auth()
	return nil
}

func (c *OpenttdAdminConnection) Auth() (err error) {
	log.Info("Authenticating with ", c.Hostname)
	out := bytes.Buffer{}
	out.WriteString("nu2ana")
	out.WriteByte(byte(0))
	out.WriteString("gopenttd")
	out.WriteByte(byte(0))
	out.WriteString("test")
	out.WriteByte(byte(0))
	c.Send(openttd_consts_admin.AdminJoin, &out)
	return
}

func (c *OpenttdAdminConnection) Close() (err error) {
	log.Debugf("Connection to %s:%d closed", c.Hostname, c.Port)
	return c.Connection.Close()
}

func (c *OpenttdAdminConnection) Send(packetType uint8, buffer *bytes.Buffer) (err error) {
	// Set a timeout
	var timeout = time.Second * 10
	c.Connection.SetDeadline(time.Now().Add(timeout))

	// Build the packet
	msg := new(bytes.Buffer)
	msgLength := uint32(buffer.Len() + 3)
	binary.Write(msg, binary.LittleEndian, msgLength)

	msg.WriteByte(packetType)
	msg.Write(buffer.Bytes())

	sendLen, err := c.Connection.Write(msg.Bytes())
	if err != nil {
		return err
	}
	log.Debug("Sent ", sendLen, " bytes")
	return nil
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
// Connections time out after 10 seconds
func ScanServerAdm(host string, port int) (err error) {
	obj := OpenttdAdminConnection{Hostname: host, Port: port}
	err = obj.Open()
	if err != nil {
		return err
	}
	defer obj.Close()
	return
}
