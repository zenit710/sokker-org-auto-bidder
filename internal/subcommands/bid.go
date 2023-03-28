package subcommands

import (
	"errors"
	"fmt"
	"log"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
)

var _ Subcommand = &bidSubcommand{}

type bidSubcommand struct {
	r player.PlayerRepository
	c client.Client
}

func NewBidSubcommand(r player.PlayerRepository, c client.Client) *bidSubcommand {
	return &bidSubcommand{r: r, c: c}
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

	// bid players
	for _, player := range players {
		err := s.handlePlayer(player)
		if err != nil {
			fmt.Printf("(%d): %v\n", player.Id, err)
		}
	}

	return nil
}

func (s *bidSubcommand) handlePlayer(p *model.Player) error {
	log.Printf("%v", p)

	info, err := s.c.FetchPlayerInfo(p.Id)
	if err != nil {
		return err
	}

	if info.Transfer.Price.MinBid.Value > p.MaxPrice {
		return errors.New("max price reached, cannot bid further")
	}

	if err = s.c.Auth(); err != nil {
		return fmt.Errorf("authorization error")
	}

	_, err = s.c.Bid(p.Id, info.Transfer.Price.MinBid.Value)
	if err != nil {
		return fmt.Errorf("bid could not be made: %v", err)
	}

	return nil
}
