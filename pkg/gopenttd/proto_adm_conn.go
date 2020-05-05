package gopenttd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
)

type OpenttdAdminConnection struct {
	Hostname   string
	Port       int
	ClientName string
	Password   string
	Connection *net.TCPConn
	Reader     *PacketReader
	Writer     *PacketWriter

	handlers      []chan packets.AdminResponsePacket
	handlersMutex sync.RWMutex

	// Set if the server has an active connection and has authenticated successfully.
	connected bool
}

func (c *OpenttdAdminConnection) Subscribe(ch chan packets.AdminResponsePacket) {
	c.handlersMutex.Lock()
	c.handlers = append(c.handlers, ch)
	log.Debugf("Subscribing handler")
	c.handlersMutex.Unlock()
}

func (c *OpenttdAdminConnection) Unsubscribe(ch chan packets.AdminResponsePacket) {
	c.handlersMutex.Lock()
	if c.handlers != nil {
		for i := range c.handlers {
			if c.handlers[i] == ch {
				c.handlers = append(c.handlers[:i], c.handlers[i+1:]...)
				break
			}
		}
	}
	log.Debugf("Unsubscribing handler")
	c.handlersMutex.Unlock()
}

func (c *OpenttdAdminConnection) publish(packet packets.AdminResponsePacket) {
	c.handlersMutex.RLock()
	// this is done because the slices refer to same array even though they are passed by value
	// thus we are creating a new slice with our elements thus preserve locking correctly.
	channels := append([]chan packets.AdminResponsePacket{}, c.handlers...)
	go func(packet packets.AdminResponsePacket, channels []chan packets.AdminResponsePacket) {
		for _, ch := range channels {
			ch <- packet
		}
	}(packet, channels)
	c.handlersMutex.RUnlock()
}

type PacketReader struct {
	reader *bufio.Reader
}

func NewPacketReader(reader io.Reader) *PacketReader {
	return &PacketReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *PacketReader) Read() (packet packets.AdminResponsePacket, err error) {
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

	return admin.Unpack(data)
}

type PacketWriter struct {
	writer io.Writer
}

func NewPacketWriter(writer io.Writer) *PacketWriter {
	return &PacketWriter{
		writer: writer,
	}
}

func (w *PacketWriter) Write(packet packets.AdminRequestPacket) (err error) {
	// Build the packet
	data := packet.Pack()
	msg := new(bytes.Buffer)
	// Length is +3 because of the metadata fields we're adding at the beginning
	msgLength := uint16(data.Len() + 3)
	binary.Write(msg, binary.LittleEndian, msgLength)

	packetType, err := admin.GetRequestPacketType(packet)
	if err != nil {
		return err
	}
	msg.WriteByte(packetType)
	msg.Write(data.Bytes())

	sendLen, err := w.writer.Write(msg.Bytes())
	if err != nil {
		return err
	}
	log.Debug("Sent ", sendLen, " bytes")
	return nil
}
