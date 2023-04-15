package subcommands

import (
	"flag"
	"fmt"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"sokker-org-auto-bidder/tools"
	"time"

	log "github.com/sirupsen/logrus"
)

var _ Subcommand = &playerAddSubcommand{}

// playerAddSubcommand handle player add command
type playerAddSubcommand struct {
	c client.Client
	r player.PlayerRepository
	fs *flag.FlagSet

	playerId uint
	maxPrice uint
}

// NewPlayerAddSubcommand returns new subcommand for adding players to the DB
func NewPlayerAddSubcommand(r player.PlayerRepository, c client.Client) *playerAddSubcommand {
	log.Trace("creating new player add subcommand handler")
	cmd := &playerAddSubcommand{
		c: c,
		r: r,
		fs: flag.NewFlagSet("add", flag.ExitOnError),
	}

	log.Trace("register command flags")
	cmd.fs.UintVar(&cmd.playerId, "playerId", 0, "Player ID")
	cmd.fs.UintVar(&cmd.maxPrice, "maxPrice", 0, "Maxium price for player to bid")

	return cmd
}

// Init parses args before run
func (s *playerAddSubcommand) Init(args []string) error {
	log.Trace("parse command flags")
	return s.fs.Parse(args)
}

// Run executes command and add player to the bid list eventually
func (s *playerAddSubcommand) Run() error {
	log.Tracef("execute player (%d) add subcommand", s.playerId)

	log.Tracef("validate player (%d) model before api calls", s.playerId)
	player := &model.Player{
		Id: s.playerId,
		MaxPrice: s.maxPrice,
	}
	if err := player.Validate(); err != nil {
		return err
	}

	log.Debugf("fetch info about player (%d)", s.playerId)
	info, err := s.c.FetchPlayerInfo(s.playerId)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not fetch player (%d) transfer details", s.playerId)
	}

	log.Tracef("check can make player (%d) bid (max price vs current price)", s.playerId)
	if s.maxPrice < info.Transfer.Price.MinBid.Value {
		return fmt.Errorf("minimum price for player (%d) is %d", s.playerId, info.Transfer.Price.MinBid.Value)
	}

	log.Tracef("parse player (%d) transfer deadline time", s.playerId)
	dt, err := tools.TimeInZone(client.TimeLayout, info.Transfer.Deadline.Date, info.Transfer.Deadline.Timezone)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not parse player (%d) transfer deadlin time", s.playerId)
	}

	log.Tracef("set player (%d) transfer deadline in player model", s.playerId)
	player.Deadline = dt.In(time.UTC)

	log.Tracef("validate player (%d) model before save", s.playerId)
	if err := player.Validate(); err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("add player (%d) to the bid list", s.playerId)
	if err := s.r.Add(player); err != nil {
		log.Error(err)
		return fmt.Errorf("player (%d) could not be added to the bid list", s.playerId)
	}

	fmt.Printf("player (%d) added to bid list: %v\n", s.playerId, player)

	return nil
}
