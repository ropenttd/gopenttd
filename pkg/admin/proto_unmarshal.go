package admin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// this is HORRIBLY poorly written
// ... but it works
// please feel free to PR improve this if you can understand how the fuck I wrote it
// (it was midnight on a saturday don't judge)

func ottdUnmarshal(data []byte, p interface{}) (err error) {
	buffer := bytes.NewBuffer(data)
	// Some funky bullshit that iterates through the elements in the ResponsePacket and unpacks them from the buffer
	// based on their type
	// I'm not even going to pretend to know how this works
	v := reflect.ValueOf(p).Elem()
	if !v.IsValid() {
		return errors.New(fmt.Sprint(v, "is invalid"))
	}
	for i := 0; i < v.NumField(); i++ {
		ottdUnmarshalData(buffer, v.Field(i), true)
	}
	return err
}

func ottdUnmarshalData(buffer *bytes.Buffer, val reflect.Value, set bool) (gen interface{}) {
	switch val.Kind() {
	// binary.Read() doesn't appear to work here (always returns 0?) so do things the long way
	// i.e binary.Read(buffer, binary.LittleEndian, nv)
	case reflect.Bool:
		var nv bool
		nv = uint8(buffer.Next(1)[0]) != 0
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.Uint8:
		var nv uint8
		nv = uint8(buffer.Next(1)[0])
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.Uint16:
		var nv uint16
		nv = binary.LittleEndian.Uint16(buffer.Next(2))
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.Uint32:
		var nv uint32
		nv = binary.LittleEndian.Uint32(buffer.Next(4))
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.Int64:
		var nv int64
		binary.Read(buffer, binary.LittleEndian, &nv)
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.Uint64:
		var nv uint64
		nv = binary.LittleEndian.Uint64(buffer.Next(8))
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.String:
		nvBytes, _ := buffer.ReadBytes(byte(0))
		nv := string(bytes.Trim(nvBytes, "\x00"))
		if set {
			val.Set(reflect.ValueOf(nv))
		}
		return nv
	case reflect.Map:
		// This is handled in a special, openttd way
		val.Set(reflect.MakeMap(val.Type()))
		var next bool
		next = uint8(buffer.Next(1)[0]) != 0
		for next {
			// there are settings to read
			// read the key data
			k := reflect.Zero(val.Type().Key())
			k = reflect.ValueOf(ottdUnmarshalData(buffer, k, false))

			// then read the value data
			v := reflect.Zero(val.Type().Elem())
			v = reflect.ValueOf(ottdUnmarshalData(buffer, v, false))

			val.SetMapIndex(k, v)

			next = uint8(buffer.Next(1)[0]) != 0
		}
		return val
	}
	return nil
}
