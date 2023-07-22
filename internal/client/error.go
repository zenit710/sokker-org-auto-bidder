package client

import (
	"errors"
	"fmt"
)

var (
	// ErrBadCredentials is raised when user login/pass is not correct
	ErrBadCredentials = errors.New("bad credentials")
	// ErrUnauthorized is raised when trying to get resources without authentication
	ErrUnauthorized = errors.New("unauthorized")
)

// ErrRequestFailed is raised when http request can not be send
type ErrRequestFailed struct {
	Type string
}

func (e *ErrRequestFailed) Error() string {
	return fmt.Sprintf("could not send '%s' request", e.Type)
}

// ErrResponseParseFailed is raised when http response can not be parsed
type ErrResponseParseFailed struct {
	Type string
}

func (e *ErrResponseParseFailed) Error() string {
	return fmt.Sprintf("'%s' response could not be parsed", e.Type)
}

// ErrResourceUnavailable is raised when can not fetch resource
type ErrResourceUnavailable struct {
	resource string
}

func (e *ErrResourceUnavailable) Error() string {
	return fmt.Sprintf("resouce unavailable: %s", e.resource)
}

// ErrNoFunds is raised when can not bid player because lack of funds
type ErrNoFunds struct {
	id uint
}

func (e *ErrNoFunds) Error() string {
	return fmt.Sprintf("no funds for player (%d) bid", e.id)
}
