package admin

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/internal/openttd_packets_admin"
	"github.com/ropenttd/gopenttd/pkg/admin/packets"
	"github.com/ropenttd/gopenttd/pkg/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type OpenttdAdminConnection struct {
	Hostname   string
	Port       int
	ClientName string
	Password   string
	Connection *net.TCPConn
	Reader     *packetReader
	Writer     *packetWriter

	callbacks CallbackMux

	// Set if the server has an active connection and has authenticated successfully.
	connected bool
}

func NewAdminConnection(hostname string, port int, password string, clientname string) (connection OpenttdAdminConnection) {
	connection.Hostname = hostname
	connection.Port = port
	connection.Password = password
	connection.ClientName = clientname
	return
}

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
	c.Reader = newPacketReader(c.Connection)
	c.Writer = newPacketWriter(c.Connection)

	log.Info("Admin connection open to ", server)
	err = c.auth()
	if err != nil {
		c.Close()
		log.Error(err)
	}
	return
}

func (c *OpenttdAdminConnection) auth() (err error) {
	log.Info("Authenticating with ", c.Hostname)
	pack := packets.AdminJoin{
		Password:   c.Password,
		ClientName: c.ClientName,
		Version:    "testing",
	}
	c.Writer.Write(pack)

	// Watch the first packet ourselves (yes I'm aware this is a little hacky)
	packet, err := c.readPacket()

	if err == io.EOF {
		// Some kind of connection error
		return util.ErrAuthentication
	} else if err != nil {
		log.Error("Read error during connection: ", err)
		return err
	}

	switch packet.(type) {
	case *packets.ServerProtocol:
		log.Info("Authenticated successfully")
		c.connected = true
	case *packets.ServerFull:
		return errors.New("server reported authentication failure - server full")
	case *packets.ServerBanned:
		return errors.New("server reported authentication failure - banned")
	case *packets.ServerError:
		return errors.New("server reported authentication failure - server error")
	case *packets.ServerShutdown:
		return errors.New("server reported authentication failure - server shutting down")
	default:
		return errors.New("unexpected initial packet")
	}

	return
}

func (c *OpenttdAdminConnection) Close() (err error) {
	c.Writer.Write(packets.AdminQuit{})
	log.Infof("Connection to %s:%d closed", c.Hostname, c.Port)
	c.connected = false
	return c.Connection.Close()
}

// readPacket is a non-public packet reader which does not check if the connection is open.
func (c *OpenttdAdminConnection) readPacket() (packet packets.AdminResponsePacket, err error) {
	data, err := c.Reader.Read()

	if err != nil {
		// Some kind of connection error
		log.Error("Read error: ", err)
		return nil, err
	}

	packet, err = c.handlePacket(data)

	if err != nil {
		// Some kind of connection error
		log.Error("Unmarshal error: ", err)
		return nil, err
	}

	return packet, err
}

func (c *OpenttdAdminConnection) ReadPacket() (packet packets.AdminResponsePacket, err error) {
	if !c.connected {
		return nil, util.ErrNotConnected
	}
	return c.readPacket()
}

type packetReader struct {
	reader *bufio.Reader
}

func newPacketReader(reader io.Reader) *packetReader {
	return &packetReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *packetReader) Read() (packet []byte, err error) {
	// Read the first part
	lengthBytes := make([]byte, 2)
	_, err = r.reader.Read(lengthBytes)
	if err != nil {
		return nil, err
	}
	if len(lengthBytes) < 2 {
		return nil, errors.New("received a packet shorter than the required length")
	}

	packLength := int(binary.LittleEndian.Uint16(lengthBytes))

	data := make([]byte, packLength-2)

	readLen, err := r.reader.Read(data)
	if err != nil {
		return nil, err
	}

	log.Debug("Response hexadecimal: ", hex.EncodeToString(data))

	if readLen+2 != packLength {
		log.Error("invalid reported buffer length: got ", readLen, ", expected ", packLength)
		// ignore the packet
		return nil, errors.New(fmt.Sprint("invalid reported buffer length: got ", readLen, ", expected ", packLength))
	}

	return data, nil
}

type packetWriter struct {
	writer io.Writer
}

func newPacketWriter(writer io.Writer) *packetWriter {
	return &packetWriter{
		writer: writer,
	}
}

func (w *packetWriter) Write(packet packets.AdminRequestPacket) (err error) {
	// Build the packet
	data := packet.Pack()
	msg := new(bytes.Buffer)
	// Length is +3 because of the metadata fields we're adding at the beginning
	msgLength := uint16(data.Len() + 3)
	binary.Write(msg, binary.LittleEndian, msgLength)

	packetType, err := openttd_packets_admin.GetRequestPacketType(packet)
	if err != nil {
		return err
	}
	msg.WriteByte(uint8(packetType))
	msg.Write(data.Bytes())

	sendLen, err := w.writer.Write(msg.Bytes())
	if err != nil {
		return err
	}
	log.Debug("Sent ", sendLen, " bytes")
	return nil
}
