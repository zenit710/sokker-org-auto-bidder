package subcommands

import (
	"fmt"
	"sokker-org-auto-bidder/internal/client"
)

var _ Subcommand = &checkAuthSubcommand{}

type checkAuthSubcommand struct {
	c client.Client
}

func NewCheckAuthSubcommand(c client.Client) *checkAuthSubcommand {
	return &checkAuthSubcommand{c: c}
}

func (s *checkAuthSubcommand) Init(args []string) error {
	return nil
}

func (s *checkAuthSubcommand) Run() error {
	club, err := s.c.Auth()
	if err != nil {
		return err
	}

	fmt.Printf("Auth success! Club ID: %d\n", club.Team.Id)
	
	return nil
}
