package errors

import "errors"

// Session errors
var (
	HaveNoBussinessId = errors.New("have no business id")
	UserAlredyLogedIn = errors.New("user alredy loged in")
)

// Ws Cmd errors
var (
	InvalidCommand       = errors.New("invalid command")
	InvalidBodyForCmd    = errors.New("invalid request body for cmd")
	InvalidTypeOfTextMsg = errors.New("invalid type of text message")
)

var InvalidTypeOfResponse = errors.New("invalid type of response")
