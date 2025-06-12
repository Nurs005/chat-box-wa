package errors

import "errors"

// Session errors
var (
	HaveNoBusinessId  = errors.New("have no business id")
	UserAlreadyLogged = errors.New("user already logged in")
)

// Ws Cmd errors
var (
	InvalidCommand       = errors.New("invalid command")
	InvalidBodyForCmd    = errors.New("invalid request body for cmd")
	InvalidTypeOfTextMsg = errors.New("invalid type of text message")
)

var InvalidTypeOfResponse = errors.New("invalid type of response")
