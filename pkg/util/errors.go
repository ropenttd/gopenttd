package util

import "errors"

var ErrNotConnected = errors.New("not connected to server")
var ErrAuthentication = errors.New("authentication failed")
