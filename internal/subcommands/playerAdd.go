package subcommands

import (
	"flag"
	"fmt"
	"log"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
)

var _ Subcommand = &playerAddSubcommand{}

type playerAddSubcommand struct {
	c client.Client
	r player.PlayerRepository
	fs *flag.FlagSet

	playerId uint
	maxPrice uint
}

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

func (s *playerAddSubcommand) Init(args []string) error {
	return s.fs.Parse(args)
}

func (s *playerAddSubcommand) Run() error {
	info, err := s.c.FetchPlayerInfo(s.playerId)
	if err != nil {
		return err
	}

	if s.maxPrice < info.Transfer.Price.MinBid.Value {
		return fmt.Errorf("minimum price for this player is %d", info.Transfer.Price.MinBid.Value)
	}

	// create player model
	player := &model.Player{
		Id: s.playerId,
		MaxPrice: s.maxPrice,
		Deadline: info.Transfer.Deadline.Date,
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
