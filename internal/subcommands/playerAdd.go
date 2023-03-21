package subcommands

import (
	"flag"
	"log"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
)

var _ Subcommand = &PlayerAddSubcommand{}

type PlayerAddSubcommand struct {
	R player.PlayerRepository
	Args []string
}

func (s *PlayerAddSubcommand) Run() error {
	// define 'add' subcommand flags set
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	playerId := addCmd.Uint("playerId", 0, "Player ID")
	maxPrice := addCmd.Uint("maxPrice", 0, "Maxium price for player to bid")

	// parse cmd flags
	addCmd.Parse(s.Args)

	// create player model
	player := &model.Player{
		Id: uint(*playerId),
		MaxPrice: uint(*maxPrice),
	}

	// validate player model
	if err := player.Validate(); err != nil {
		return err
	}

	// save player into the DB
	if err := s.R.Add(player); err != nil {
		return err
	}

	log.Printf("player added to bid list: %v", player)

	return nil
}
