package gopenttd

import "errors"

var errNotConnected = errors.New("not connected to server")
var errAuthentication = errors.New("authentication failed")
