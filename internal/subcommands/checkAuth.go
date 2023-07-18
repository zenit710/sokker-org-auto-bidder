package subcommands

import (
	"errors"
	"fmt"
	"sokker-org-auto-bidder/internal/client"

	log "github.com/sirupsen/logrus"
)

var _ Subcommand = &checkAuthSubcommand{}

// checkAuthSubcommand handle auth check command
type checkAuthSubcommand struct {
	c client.Client
}

// NewCheckAuthSubcommand returns new subcommand for checking auth
func NewCheckAuthSubcommand(c client.Client) *checkAuthSubcommand {
	log.Trace("creating new check auth subcommand handler")
	return &checkAuthSubcommand{c: c}
}

// Init parse command run args
func (s *checkAuthSubcommand) Init(args []string) error {
	return nil
}

// Run executes subcommand
func (s *checkAuthSubcommand) Run() (interface{}, error) {
	log.Trace("execute check auth subcommand")

	log.Debug("auth in sokker.org")
	club, err := s.c.Auth()
	if err != nil && !errors.Is(err, client.ErrBadCredentials) {
		log.Error(err)
		return nil, fmt.Errorf("authorization error")
	}

	if club == nil {
		fmt.Printf("Auth failed.\n")
	} else {
		fmt.Printf("Auth success! Club ID: %d\n", club.Team.Id)
	}

	return nil, nil
}
