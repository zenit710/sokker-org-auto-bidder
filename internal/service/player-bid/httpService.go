package playerbid

import (
	"fmt"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/tools"
	"time"

	log "github.com/sirupsen/logrus"
)

var _ PlayerBidService = &httpPlayerBidService{}

// httpPlayerBidService handles player bid process
type httpPlayerBidService struct {
	r player.PlayerRepository
	c client.Client
}

// NewHttpPlayerBidService crates new httpPlayerBidService instance
func NewHttpPlayerBidService(r player.PlayerRepository, c client.Client) *httpPlayerBidService {
	return &httpPlayerBidService{r, c}
}

// Bid handle player bid process
func (s *httpPlayerBidService) Bid(p *model.Player, clubId uint) error {
	log.Tracef("handle player (%d) bid", p.Id)

	log.Debugf("fetch player (%d) transfer info", p.Id)
	info, err := s.c.FetchPlayerInfo(p.Id)
	if err != nil {
		log.Error(err)
		return &ErrTransferInfoFetchFailed{p.Id}
	}

	log.Tracef("check can player (%d) bid be made (value vs. max price)", p.Id)
	if info.Transfer.Price.MinBid.Value > p.MaxPrice {
		if err = s.r.Delete(p); err != nil {
			log.Error(err)
			fmt.Printf("player (%d) did not remove from bid list, something went wrong\n", p.Id)
		}

		return &ErrMaxPriceReached{p.Id}
	}

	log.Tracef("check is player (%d) bid neccessary (current leader)", p.Id)
	if info.Transfer.BuyerId == clubId {
		return &ErrCurrentLeader{p.Id}
	}

	log.Debugf("make player (%d) bid", p.Id)
	tr, err := s.c.Bid(p.Id, info.Transfer.Price.MinBid.Value)
	if err != nil {
		log.Error(err)
		return &ErrCouldNotBid{p.Id, err}
	}

	log.Tracef("parse player (%d) transfer deadline time", p.Id)
	newDeadline, err := tools.TimeInZone(client.TimeLayout, tr.Deadline.Date.Date, tr.Deadline.Date.Timezone)
	if err != nil {
		log.Error(err)
		return &ErrDeadlineParse{p.Id}
	}

	log.Tracef("check is player (%d) transfer deadline still valid", p.Id)
	if p.Deadline.Before(newDeadline) {
		p.Deadline = newDeadline.In(time.UTC)

		log.Debugf("update player (%d) transfer deadline", p.Id)
		if err = s.r.Update(p); err != nil {
			log.Error(err)
			return &ErrDeadlineNotUpdated{p.Id}
		}
	}

	return nil
}
