package subcommands

import (
	"errors"
	"fmt"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/tools"
	"time"

	log "github.com/sirupsen/logrus"
)

var _ Subcommand = &bidSubcommand{}

// bidSubcommand handle player bid action
type bidSubcommand struct {
	r player.PlayerRepository
	c client.Client
}

// NewBidSubcommand returns new player bid command handler
func NewBidSubcommand(r player.PlayerRepository, c client.Client) *bidSubcommand {
	log.Trace("creating new bid subcommand handler")
	return &bidSubcommand{r: r, c: c}
}

// Init does nothing, bid command do not require extra arguments
func (s *bidSubcommand) Init(args []string) error {
	return nil
}

// Run executes bid subcommand
func (s *bidSubcommand) Run() error {
	log.Trace("make bid for listed players")

	log.Trace("fetch players to bid")
	players, err := s.r.List()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("auth in sokker.org")
	club, err := s.c.Auth()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("authorization error")
	}

	log.Trace("make players bids")
	for _, player := range players {
		err := s.handlePlayer(player, club.Team.Id)
		if err != nil {
			fmt.Printf("(%d): %v\n", player.Id, err)
		} else {
			fmt.Printf("(%d): bid made\n", player.Id)
		}
	}

	return nil
}

// handlePlayer handle player bid process
func (s *bidSubcommand) handlePlayer(p *model.Player, clubId uint) error {
	log.Tracef("handle player %d bid", p.Id)

	log.Tracef("fetch transfer info")
	info, err := s.c.FetchPlayerInfo(p.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("check can bid be made (value vs. max price)")
	if info.Transfer.Price.MinBid.Value > p.MaxPrice {
		if err = s.r.Delete(p); err != nil {
			log.Error(err)
			fmt.Println("player did not remove from bid list, something went wrong")
		}

		return errors.New("max price reached, cannot bid further")
	}

	log.Trace("check is bid neccessary (current leader)")
	if info.Transfer.BuyerId == clubId {
		return errors.New("you are current leader, no reason to bid")
	}

	log.Trace("make player bid")
	tr, err := s.c.Bid(p.Id, info.Transfer.Price.MinBid.Value)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("bid could not be made: %v", err)
	}

	log.Trace("parse transfer deadline time")
	newDeadline, err := tools.TimeInZone(client.TimeLayout, tr.Deadline.Date.Date, tr.Deadline.Date.Timezone)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("check is transfer deadline still valid")
	if p.Deadline.Before(newDeadline) {
		p.Deadline = newDeadline.In(time.UTC)

		log.Trace("update player transfer deadline")
		if err = s.r.Update(p); err != nil {
			fmt.Println("player transfer deadline was not updated, it can lead to mistakes, sorry")
		}
	}

	return nil
}
