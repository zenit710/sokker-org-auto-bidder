package client

import (
	"errors"
)

var (
	// ErrBadCredentials is raised when user login/pass is not correct
	ErrBadCredentials = errors.New("bad credentials")
)
