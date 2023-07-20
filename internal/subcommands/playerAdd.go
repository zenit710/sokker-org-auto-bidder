package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/tools"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	_                  Subcommand = &playerAddSubcommand{}
	requiredFlags                 = []string{"playerId", "maxPrice"}
	ErrCanNotParseArgs            = errors.New("could not parse subcommand args")
)

// playerAddSubcommand handle player add command
type playerAddSubcommand struct {
	c  client.Client
	r  player.PlayerRepository
	fs *flag.FlagSet

	playerId uint
	maxPrice uint
}

// NewPlayerAddSubcommand returns new subcommand for adding players to the DB
func NewPlayerAddSubcommand(r player.PlayerRepository, c client.Client) *playerAddSubcommand {
	log.Trace("creating new player add subcommand handler")
	cmd := &playerAddSubcommand{
		c:  c,
		r:  r,
		fs: flag.NewFlagSet("add", flag.ContinueOnError),
	}

	log.Trace("register command flags")
	cmd.fs.UintVar(&cmd.playerId, "playerId", 0, "Player ID")
	cmd.fs.UintVar(&cmd.maxPrice, "maxPrice", 0, "Maxium price for player to bid")

	return cmd
}

// Init parses args before run
func (s *playerAddSubcommand) Init(args []string) error {
	log.Trace("parse add subcommand args")
	if err := s.fs.Parse(args); err != nil {
		log.Error(err)
		return ErrCanNotParseArgs
	}

	log.Trace("verify required flags")
	reqFlagsProvided := []string{}
	s.fs.Visit(func(f *flag.Flag) {
		if stringSliceContains(requiredFlags, f.Name) {
			reqFlagsProvided = append(reqFlagsProvided, f.Name)
		}
	})
	if len(reqFlagsProvided) < len(requiredFlags) {
		return &ErrMissingFlags{requiredFlags}
	}

	return nil
}

// Run executes command and add player to the bid list eventually
func (s *playerAddSubcommand) Run() (interface{}, error) {
	log.Tracef("execute player (%d) add subcommand", s.playerId)

	log.Debugf("fetch info about player (%d)", s.playerId)
	info, err := s.c.FetchPlayerInfo(s.playerId)
	if err != nil {
		log.Error(err)
		return nil, &ErrCanNotFetchTransferDetails{}
	}

	log.Tracef("check can make player (%d) bid (max price vs current price)", s.playerId)
	if s.maxPrice < info.Transfer.Price.MinBid.Value {
		return nil, &ErrMaxPriceToLow{}
	}

	log.Tracef("parse player (%d) transfer deadline time", s.playerId)
	dt, err := tools.TimeInZone(client.TimeLayout, info.Transfer.Deadline.Date, info.Transfer.Deadline.Timezone)
	if err != nil {
		log.Error(err)
		return nil, &ErrCanNotParseTransferDeadline{}
	}

	log.Tracef("map player (%d) from response to player model", s.playerId)
	player := &model.Player{
		Id:       s.playerId,
		MaxPrice: s.maxPrice,
		Deadline: dt.In(time.UTC),
	}

	log.Tracef("validate player (%d) model before save", s.playerId)
	if err := player.Validate(); err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debugf("add player (%d) to the bid list", s.playerId)
	if err := s.r.Add(player); err != nil {
		log.Error(err)
		return nil, &ErrCanNotAddToBidList{}
	}

	fmt.Printf("player (%d) added to bid list: %v\n", s.playerId, player)

	return nil, nil
}

func stringSliceContains(slice []string, search string) bool {
	for _, s := range slice {
		if s == search {
			return true
		}
	}
	return false
}

type ErrCanNotFetchTransferDetails struct {
	id uint
}

func (e *ErrCanNotFetchTransferDetails) Error() string {
	return fmt.Sprintf("could not fetch player (%d) transfer details", e.id)
}

type ErrMaxPriceToLow struct {
	id, minPrice uint
}

func (e *ErrMaxPriceToLow) Error() string {
	return fmt.Sprintf("minimum price for player (%d) is %d", e.id, e.minPrice)
}

type ErrCanNotParseTransferDeadline struct {
	id uint
}

func (e *ErrCanNotParseTransferDeadline) Error() string {
	return fmt.Sprintf("could not parse player (%d) transfer deadline time", e.id)
}

type ErrCanNotAddToBidList struct {
	id uint
}

func (e *ErrCanNotAddToBidList) Error() string {
	return fmt.Sprintf("player (%d) could not be added to the bid list", e.id)
}
