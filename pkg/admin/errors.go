package admin

import "errors"

var ErrAlreadyConnected = errors.New("admin connection already open")

var ErrInvalidUpdateFrequency = errors.New("given update frequency is not valid")

// ErrNilState is returned when the state is nil.
var ErrNilState = errors.New("state not instantiated, please use admin.New() or assign Session.State")

// ErrStateNotFound is returned when the state cache
// requested is not found
var ErrStateNotFound = errors.New("state cache not found")

var ErrServerFull = errors.New("server is full")
var ErrServerBanned = errors.New("banned from server")
var ErrServerError = errors.New("server encountered an error")
