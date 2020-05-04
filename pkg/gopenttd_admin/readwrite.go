package gopenttd_admin

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/internal/openttd_consts_admin"
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
	Reader     *PacketReader
	Writer     *PacketWriter
}

type PacketReader struct {
	reader *bufio.Reader
}

func NewPacketReader(reader io.Reader) *PacketReader {
	return &PacketReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *PacketReader) Read() (packet openttd_consts_admin.AdminPacket, err error) {
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

	return openttd_consts_admin.Unpack(data)
}

type PacketWriter struct {
	writer io.Writer
}

func NewPacketWriter(writer io.Writer) *PacketWriter {
	return &PacketWriter{
		writer: writer,
	}
}

func (w *PacketWriter) Write(packet openttd_consts_admin.AdminPacket) (err error) {
	// Build the packet
	data := packet.Data()
	msg := new(bytes.Buffer)
	// Length is +3 because of the metadata fields we're adding at the beginning
	msgLength := uint16(data.Len() + 3)
	binary.Write(msg, binary.LittleEndian, msgLength)

	packetType, err := openttd_consts_admin.GetAdminPacketType(packet)
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
