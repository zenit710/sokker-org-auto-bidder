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
	log.Tracef("execute player %d add subcommand", s.playerId)

	log.Trace("fetch info about player")
	info, err := s.c.FetchPlayerInfo(s.playerId)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("check can make bid (max price vs current price)")
	if s.maxPrice < info.Transfer.Price.MinBid.Value {
		return fmt.Errorf("minimum price for this player is %d", info.Transfer.Price.MinBid.Value)
	}

	log.Trace("parse transfer deadline time")
	dt, err := tools.TimeInZone(client.TimeLayout, info.Transfer.Deadline.Date, info.Transfer.Deadline.Timezone)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace("map player info from response to player model")
	player := &model.Player{
		Id: s.playerId,
		MaxPrice: s.maxPrice,
		Deadline: dt.In(time.UTC),
	}

	log.Trace("validate player model before save")
	if err := player.Validate(); err != nil {
		return err
	}

	log.Trace("add player to the bid list")
	if err := s.r.Add(player); err != nil {
		return err
	}

	fmt.Printf("player added to bid list: %v\n", player)

	return nil
}
