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
	// there's probably a better way to do this but oh well this is quite performant
	case iServerFull:
		obj := packets.ServerFull{}
		packet = obj
	case iServerBanned:
		obj := packets.ServerBanned{}
		packet = obj
	case iServerError:
		obj := packets.ServerError{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerProtocol:
		obj := packets.ServerProtocol{}
		err = obj.Unpack(reader)
		packet = obj
	case iServerWelcome:
		obj := packets.ServerWelcome{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerNewgame:
		obj := packets.ServerNewgame{}
		packet = obj
	case iServerShutdown:
		obj := packets.ServerShutdown{}
		packet = obj
	case iServerDate:
		obj := packets.ServerDate{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerClientJoin:
		obj := packets.ServerClientJoin{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerClientInfo:
		obj := packets.ServerClientInfo{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerClientUpdate:
		obj := packets.ServerClientUpdate{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerClientQuit:
		obj := packets.ServerClientQuit{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerClientError:
		obj := packets.ServerClientError{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCompanyNew:
		obj := packets.ServerCompanyNew{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCompanyInfo:
		obj := packets.ServerCompanyInfo{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCompanyUpdate:
		obj := packets.ServerCompanyUpdate{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCompanyRemove:
		obj := packets.ServerCompanyRemove{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCompanyEconomy:
		obj := packets.ServerCompanyEconomy{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCompanyStats:
		obj := packets.ServerCompanyStats{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerChat:
		obj := packets.ServerChat{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerRcon:
		obj := packets.ServerRcon{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerCmdNames:
		obj := packets.ServerCmdNames{}
		err = obj.Unpack(reader)
		packet = obj
	case iServerCmdLogging:
		obj := packets.ServerCmdLogging{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerGamescript:
		obj := packets.ServerGamescript{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerRconEnd:
		obj := packets.ServerRconEnd{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerConsole:
		obj := packets.ServerConsole{}
		err = genericUnpack(&obj, reader)
		packet = obj
	case iServerPong:
		obj := packets.ServerPong{}
		err = genericUnpack(&obj, reader)
		packet = obj
	}
	return packet, err
}
