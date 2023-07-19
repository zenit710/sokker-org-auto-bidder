package subcommands

import (
	"errors"
	"fmt"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/repository/player"
	playerbid "sokker-org-auto-bidder/internal/service/player-bid"

	log "github.com/sirupsen/logrus"
)

var (
	_                       Subcommand = &bidSubcommand{}
	ErrDbFetchPlayersFailed            = errors.New("could not fetch players to bid")
	ErrApiAuthFailed                   = errors.New("authorization error")
)

// bidSubcommand handle player bid action
type bidSubcommand struct {
	r player.PlayerRepository
	c client.Client
	b playerbid.PlayerBidService
}

// NewBidSubcommand returns new player bid command handler
func NewBidSubcommand(r player.PlayerRepository, c client.Client, b playerbid.PlayerBidService) *bidSubcommand {
	log.Trace("creating new bid subcommand handler")
	return &bidSubcommand{r, c, b}
}

// Init does nothing, bid command do not require extra arguments
func (s *bidSubcommand) Init(args []string) error {
	return nil
}

// Run executes bid subcommand
func (s *bidSubcommand) Run() (interface{}, error) {
	log.Trace("make bid for listed players")
	output := BidSubcommandOutput{}

	log.Debug("fetch players to bid")
	players, err := s.r.List()
	if err != nil {
		log.Error(err)
		return output, ErrDbFetchPlayersFailed
	}
	log.Debugf("%d players for bid", len(players))

	log.Debug("auth in sokker.org")
	club, err := s.c.Auth()
	if err != nil {
		log.Error(err)
		return output, ErrApiAuthFailed
	}

	log.Debug("make players bids")
	for _, player := range players {
		err := s.b.Bid(player, club.Team.Id)
		if err != nil {
			fmt.Printf("player (%d): %v\n", player.Id, err)
			output.Failed++
		} else {
			fmt.Printf("player (%d): bid made\n", player.Id)
			output.Ok++
		}
	}

	return output, nil
}

type BidSubcommandOutput struct {
	Ok, Failed uint
}
