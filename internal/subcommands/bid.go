package subcommands

import (
	"log"
	"sokker-org-auto-bidder/internal/repository/player"
)

var _ Subcommand = &BidSubcommand{}

type BidSubcommand struct {
	R player.PlayerRepository
}

func (s *BidSubcommand) Run(args []string) error {
	log.Print("make bid for listed players:")

	// get players to bid list
	players, err := s.R.GetList()
	if err != nil {
		return err
	}

	// print all players to bid
	for _, player := range players {
		log.Printf("%v", player)
	}

	return nil
}
