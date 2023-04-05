package subcommands

import (
	"fmt"
	"sokker-org-auto-bidder/internal/client"
)

var _ Subcommand = &checkAuthSubcommand{}

// checkAuthSubcommand handle auth check command
type checkAuthSubcommand struct {
	c client.Client
}

// NewCheckAuthSubcommand returns new subcommand for checking auth
func NewCheckAuthSubcommand(c client.Client) *checkAuthSubcommand {
	return &checkAuthSubcommand{c: c}
}

// Init parse command run args
func (s *checkAuthSubcommand) Init(args []string) error {
	return nil
}

// Run executes subcommand
func (s *checkAuthSubcommand) Run() error {
	club, err := s.c.Auth()
	if err != nil {
		return err
	}

	fmt.Printf("Auth success! Club ID: %d\n", club.Team.Id)
	
	return nil
}
