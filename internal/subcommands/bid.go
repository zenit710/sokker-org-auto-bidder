package subcommands

import (
	"flag"
	"log"
	"sokker-org-auto-bidder/internal/repository/player"
)

var _ Subcommand = &bidSubcommand{}

type bidSubcommand struct {
	r player.PlayerRepository
	fs *flag.FlagSet
}

func NewBidSubcommand(r player.PlayerRepository) *bidSubcommand {
	return &bidSubcommand{r: r}
}

func (s *bidSubcommand) Init(args []string) error {
	return nil
}

func (s *bidSubcommand) Run() error {
	log.Print("make bid for listed players:")

	// get players to bid list
	players, err := s.r.GetList()
	if err != nil {
		return err
	}

	// print all players to bid
	for _, player := range players {
		log.Printf("%v", player)
	}

	return nil
}
