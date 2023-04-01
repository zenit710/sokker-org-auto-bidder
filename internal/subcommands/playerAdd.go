package subcommands

import (
	"flag"
	"fmt"
	"log"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/tools"
	"time"
)

var _ Subcommand = &playerAddSubcommand{}

type playerAddSubcommand struct {
	c client.Client
	r player.PlayerRepository
	fs *flag.FlagSet

	playerId uint
	maxPrice uint
}

// NewPlayerAddSubcommand returns new subcommand for adding players to the DB
func NewPlayerAddSubcommand(r player.PlayerRepository, c client.Client) *playerAddSubcommand {
	cmd := &playerAddSubcommand{
		c: c,
		r: r,
		fs: flag.NewFlagSet("add", flag.ExitOnError),
	}

	cmd.fs.UintVar(&cmd.playerId, "playerId", 0, "Player ID")
	cmd.fs.UintVar(&cmd.maxPrice, "maxPrice", 0, "Maxium price for player to bid")

	return cmd
}

// Init parses args before run
func (s *playerAddSubcommand) Init(args []string) error {
	return s.fs.Parse(args)
}

// Run executes command and add player to the bid list eventually
func (s *playerAddSubcommand) Run() error {
	info, err := s.c.FetchPlayerInfo(s.playerId)
	if err != nil {
		return err
	}

	// check can be any bid made
	if s.maxPrice < info.Transfer.Price.MinBid.Value {
		return fmt.Errorf("minimum price for this player is %d", info.Transfer.Price.MinBid.Value)
	}

	// get transfer deadline time including timezone
	dt, err := tools.TimeInZone(client.TimeLayout, info.Transfer.Deadline.Date, info.Transfer.Deadline.Timezone)
	if err != nil {
		return err
	}

	// create player model
	player := &model.Player{
		Id: s.playerId,
		MaxPrice: s.maxPrice,
		Deadline: dt.In(time.UTC),
	}

	// validate player model
	if err := player.Validate(); err != nil {
		return err
	}

	// save player into the DB
	if err := s.r.Add(player); err != nil {
		return err
	}

	log.Printf("player added to bid list: %v", player)

	return nil
}
