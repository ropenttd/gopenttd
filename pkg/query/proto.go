package query

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/internal/openttd_packets_udp"
	"github.com/ropenttd/gopenttd/pkg/util"
	"net"
	"time"
)

func (c *OpenttdClientConnection) Open() (err error) {
	if c.Connection != nil {
		// We're already connected
		return nil
	}
	// Determine the correct UDP address
	server := fmt.Sprintf("%s:%d", c.Hostname, c.Port) // "255.255.255.255", 10000
	serverAddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return err
	}

	// Open the connection
	c.Connection, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return err
	}
	return nil
}

func (c *OpenttdClientConnection) Close() (err error) {
	return c.Connection.Close()
}

func (c *OpenttdClientConnection) Query(packetType openttd_packets_udp.UdpPacketIndex, expect openttd_packets_udp.UdpPacketIndex) (data *bytes.Buffer, err error) {
	// this buffer _should_ be long enough? probably a better way of doing this than allocating a 2m bytearray
	inBuf := make([]byte, 2048)
	var readLen int

	// Set a timeout
	var timeout = time.Second * 10
	c.Connection.SetDeadline(time.Now().Add(timeout))

	// send the Info request
	discoverMsg := []byte{3, 0, uint8(packetType)}
	_, err = c.Connection.Write(discoverMsg)
	if err != nil {
		return nil, err
	}

	for {
		// Start a loop to read from UDP, and exit out if an error or timeout is experienced
		readLen, _, err = c.Connection.ReadFromUDP(inBuf)

		if err != nil {
			// Some kind of connection error
			return nil, err
		}

		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			// Timeout error
			return nil, err
		}

		if readLen != 0 {
			// Success!
			break
		}
	}

	buffer := bytes.NewBuffer(inBuf)
	retLength := int(binary.LittleEndian.Uint16(buffer.Next(2)))
	if readLen != retLength {
		return nil, errors.New(fmt.Sprint("invalid reported buffer length: got ", retLength, ", expected ", readLen))
	}

	retType := openttd_packets_udp.UdpPacketIndex(buffer.Next(1)[0])
	if expect != retType {
		return nil, errors.New(fmt.Sprint("unexpected response packet type: got ", retType, ", expected ", expect))
	}

	return buffer, nil
}

// ScanServer takes a hostname and port and returns an OpenttdServerState struct containing the data available from it.
// Connections time out after 10 seconds
func ScanServer(host string, port int) (serverstate util.OpenttdServerState, err error) {
	obj := OpenttdClientConnection{Hostname: host, Port: port}
	err = obj.Open()
	if err != nil {
		return util.OpenttdServerState{Status: false, Error: err}, err
	}
	defer obj.Close()

	// Let's get the initial set of data using CLIENT_FIND_SERVER
	result, err := obj.Query(openttd_packets_udp.ClientFindServer, openttd_packets_udp.ServerResponse)
	if err != nil {
		return util.OpenttdServerState{Status: false, Error: err}, err
	}
	serverstate = util.OpenttdServerState{Host: obj.Hostname}
	serverstate.PopulateServerState(result)

	// Then we get the company data using CLIENT_DETAIL_INFO
	result, err = obj.Query(openttd_packets_udp.ClientDetailInfo, openttd_packets_udp.ServerDetailInfo)
	if err != nil {
		return util.OpenttdServerState{Status: false, Error: err}, err
	}
	serverstate.PopulateCompanyState(result)
	return
}
