package ecode

import "errors"

var (
	ErrNullPointer = errors.New("null reference instance")
	ErrNoCommand   = errors.New("no command")
	ErrAPIKey      = errors.New("no OpenAI API key")
	ErrStream      = errors.New("error in response stream")
	ErrNoContent   = errors.New("no content")
	ErrBinary      = errors.New("binary file")
)
