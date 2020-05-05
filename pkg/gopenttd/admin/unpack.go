package admin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ropenttd/gopenttd/pkg/gopenttd/admin/packets"
	"reflect"
)

func genericUnpack(p packets.AdminResponsePacket, buffer *bytes.Buffer) (err error) {
	// Some funky bullshit that iterates through the elements in the ResponsePacket and unpacks them from the buffer
	// based on their type
	// I'm not even going to pretend to know how this works
	v := reflect.ValueOf(p).Elem()
	if !v.IsValid() {
		return errors.New(fmt.Sprint(v, "is invalid"))
	}
	for i := 0; i < v.NumField(); i++ {
		switch f := v.Field(i); f.Kind() {
		// binary.Read() doesn't appear to work here (always returns 0?) so do things the long way
		// i.e binary.Read(buffer, binary.LittleEndian, nv)
		case reflect.Bool:
			var nv bool
			nv = uint8(buffer.Next(1)[0]) != 0
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint8:
			var nv uint8
			nv = uint8(buffer.Next(1)[0])
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint16:
			var nv uint16
			nv = binary.LittleEndian.Uint16(buffer.Next(2))
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint32:
			var nv uint32
			nv = binary.LittleEndian.Uint32(buffer.Next(4))
			f.Set(reflect.ValueOf(nv))
		case reflect.Uint64:
			var nv uint64
			nv = binary.LittleEndian.Uint64(buffer.Next(8))
			f.Set(reflect.ValueOf(nv))
		case reflect.String:
			nvBytes, _ := buffer.ReadBytes(byte(0))
			nv := string(bytes.Trim(nvBytes, "\x00"))
			f.Set(reflect.ValueOf(nv))
		}
	}
	return err
}

func Unpack(data []byte) (packet packets.AdminResponsePacket, err error) {
	reader := bytes.NewBuffer(data)
	packetType := uint8(reader.Next(1)[0])
	switch packetType {
	case 103:
		obj := packets.ServerProtocol{Settings: map[uint16]uint16{}}
		err = obj.Unpack(reader)
		packet = obj
	case 104:
		obj := packets.ServerWelcome{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case 105:
		obj := packets.ServerNewgame{}
		packet = obj
	case 106:
		obj := packets.ServerShutdown{}
		packet = obj
	case 107:
		obj := packets.ServerDate{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case 108:
		obj := packets.ServerClientJoin{}
		err = genericUnpack(&obj, reader)
		packet = obj
	}
	return packet, err
}
