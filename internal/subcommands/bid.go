package subcommands

import (
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

	log.Debug("fetch players to bid")
	players, err := s.r.List()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not fetch players to bid")
	}
	log.Debugf("%d players for bid", len(players))

	log.Debug("auth in sokker.org")
	club, err := s.c.Auth()
	if err != nil {
		log.Error(err)
		return fmt.Errorf("authorization error")
	}

	log.Debug("make players bids")
	for _, player := range players {
		err := s.handlePlayer(player, club.Team.Id)
		if err != nil {
			fmt.Printf("player (%d): %v\n", player.Id, err)
		} else {
			fmt.Printf("player (%d): bid made\n", player.Id)
		}
	}

	return nil
}

// handlePlayer handle player bid process
func (s *bidSubcommand) handlePlayer(p *model.Player, clubId uint) error {
	log.Tracef("handle player (%d) bid", p.Id)

	log.Debugf("fetch player (%d) transfer info", p.Id)
	info, err := s.c.FetchPlayerInfo(p.Id)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not fetch player (%d) transfer info", p.Id)
	}

	log.Tracef("check can player (%d) bid be made (value vs. max price)", p.Id)
	if info.Transfer.Price.MinBid.Value > p.MaxPrice {
		if err = s.r.Delete(p); err != nil {
			log.Error(err)
			fmt.Printf("player (%d) did not remove from bid list, something went wrong\n", p.Id)
		}

		return fmt.Errorf("max price reached, cannot bid player (%d) further", p.Id)
	}

	log.Tracef("check is player (%d) bid neccessary (current leader)", p.Id)
	if info.Transfer.BuyerId == clubId {
		return fmt.Errorf("you are current leader, no reason to bid player (%d)", p.Id)
	}

	log.Debugf("make player (%d) bid", p.Id)
	tr, err := s.c.Bid(p.Id, info.Transfer.Price.MinBid.Value)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("player (%d) bid could not be made: %v", p.Id, err)
	}

	log.Tracef("parse player (%d) transfer deadline time", p.Id)
	newDeadline, err := tools.TimeInZone(client.TimeLayout, tr.Deadline.Date.Date, tr.Deadline.Date.Timezone)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not parse player (%d) transfer deadlin tine", p.Id)
	}

	log.Tracef("check is player (%d) transfer deadline still valid", p.Id)
	if p.Deadline.Before(newDeadline) {
		p.Deadline = newDeadline.In(time.UTC)

		log.Debugf("update player (%d) transfer deadline", p.Id)
		if err = s.r.Update(p); err != nil {
			log.Error(err)
			return fmt.Errorf("player (%d) transfer deadline was not updated, it can lead to mistakes, sorry", p.Id)
		}
	}

	return nil
}
