package util

import "errors"

var ErrAuthentication = errors.New("authentication failed")
var ErrBadWrite = errors.New("writing to admin port cut short")
var ErrInvalidIncomingPacket = errors.New("received an invalid packet from the server")
var ErrNotConnected = errors.New("not connected to server")
