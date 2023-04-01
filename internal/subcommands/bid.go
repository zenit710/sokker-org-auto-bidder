package subcommands

import (
	"errors"
	"fmt"
	"log"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/tools"
	"time"
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
	players, err := s.r.List()
	if err != nil {
		return err
	}

	// auth
	club, err := s.c.Auth()
	if err != nil {
		return fmt.Errorf("authorization error")
	}

	// bid players
	for _, player := range players {
		err := s.handlePlayer(player, club.Team.Id)
		if err != nil {
			fmt.Printf("(%d): %v\n", player.Id, err)
		}
	}

	return nil
}

func (s *bidSubcommand) handlePlayer(p *model.Player, clubId uint) error {
	log.Printf("%v", p)

	// fetch player info
	info, err := s.c.FetchPlayerInfo(p.Id)
	if err != nil {
		return err
	}

	// check can bid furhter
	if info.Transfer.Price.MinBid.Value > p.MaxPrice {
		if err = s.r.Delete(p); err != nil {
			fmt.Printf("player did not remove from bid list: %v", err)
		}

		return errors.New("max price reached, cannot bid further")
	}

	// check current leader
	if info.Transfer.BuyerId == clubId {
		return errors.New("you are current leader, no reason to bid")
	}

	// bid player
	tr, err := s.c.Bid(p.Id, info.Transfer.Price.MinBid.Value)
	if err != nil {
		return fmt.Errorf("bid could not be made: %v", err)
	}

	// check transfer deadline changed
	newDeadline, err := tools.TimeInZone(client.TimeLayout, tr.Deadline.Date.Date, tr.Deadline.Date.Timezone)
	if err != nil {
		return err
	}

	// update player entry with new deadline
	if p.Deadline.Before(newDeadline) {
		p.Deadline = newDeadline.In(time.UTC)
		s.r.Update(p)
	}

	return nil
}
