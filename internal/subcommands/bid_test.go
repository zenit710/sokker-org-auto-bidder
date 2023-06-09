package subcommands

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"testing"
)

func createBidSubcommand() (player.PlayerRepository, client.Client, *bidSubcommand) {
	r := player.NewMockPlayerRepository()
	c := client.NewMockClient()
	s := NewBidSubcommand(r, c)
	return r, c, s
}

func TestNewBidSubcommand(t *testing.T) {
	_, _, s := createBidSubcommand()
	if s == nil {
		t.Error("*bidSubcommand expected but nil returned")
	}
}

func TestInit(t *testing.T) {
	_, _, s := createBidSubcommand()
	if err := s.Init([]string{}); err != nil {
		t.Errorf("nil error expected, %v returned", err)
	}
}

func TestRunDbFetchFailure(t *testing.T) {
	r := player.NewMockPlayerRepository()
	r.On("List").Return([]*model.Player{}, errors.New("error"))
	c := client.NewMockClient()
	s := NewBidSubcommand(r, c)
	err := s.Run()
	if err == nil || !errors.Is(err, ErrDbFetchPlayersFailed) {
		t.Errorf("'%v' expected, '%v' returned", ErrDbFetchPlayersFailed, err)
	}
}

func TestRunApiAuthFailure(t *testing.T) {
	r := player.NewMockPlayerRepository()
	r.On("List").Return([]*model.Player{}, nil)
	c := client.NewMockClient()
	c.On("Auth").Return(c.GetEmptyClubInfoResponse(), errors.New("error"))
	s := NewBidSubcommand(r, c)
	err := s.Run()
	if err == nil || !errors.Is(err, ErrApiAuthFailed) {
		t.Errorf("'%v' expected, '%v' returned", ErrApiAuthFailed, err)
	}
}

func TestRunNoListedPlayers(t *testing.T) {
	r := player.NewMockPlayerRepository()
	r.On("List").Return([]*model.Player{}, nil)
	c := client.NewMockClient()
	c.On("Auth").Return(c.GetEmptyClubInfoResponse(), nil)
	s := NewBidSubcommand(r, c)
	err := s.Run()
	if err != nil {
		t.Errorf("nil error expected, '%v' returned", err)
	}
}

// test run with player handling E2E
// test player handle errors
