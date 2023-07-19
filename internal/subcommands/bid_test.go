package subcommands

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	playerbid "sokker-org-auto-bidder/internal/service/player-bid"
	"testing"
)

func createBidSubcommand() (*player.MockPlayerRepository, *client.MockClient, *playerbid.MockPlayerBidService, *bidSubcommand) {
	r := &player.MockPlayerRepository{}
	c := &client.MockClient{}
	b := &playerbid.MockPlayerBidService{}
	s := NewBidSubcommand(r, c, b)
	return r, c, b, s
}

func TestNewBidSubcommand(t *testing.T) {
	_, _, _, s := createBidSubcommand()
	if s == nil {
		t.Error("*bidSubcommand expected but nil returned")
	}
}

func TestInit(t *testing.T) {
	_, _, _, s := createBidSubcommand()
	if err := s.Init([]string{}); err != nil {
		t.Errorf("nil error expected, %v returned", err)
	}
}

func TestRunDbFetchFailure(t *testing.T) {
	r, _, _, s := createBidSubcommand()
	r.On("List").Return([]*model.Player{}, errors.New("error"))
	_, err := s.Run()
	if err == nil || !errors.Is(err, ErrDbFetchPlayersFailed) {
		t.Errorf("'%v' expected, '%v' returned", ErrDbFetchPlayersFailed, err)
	}
}

func TestRunApiAuthFailure(t *testing.T) {
	r, c, _, s := createBidSubcommand()
	r.On("List").Return([]*model.Player{}, nil)
	c.On("Auth").Return(c.GetClubInfoResponse(0), errors.New("error"))
	_, err := s.Run()
	if err == nil || !errors.Is(err, ErrApiAuthFailed) {
		t.Errorf("'%v' expected, '%v' returned", ErrApiAuthFailed, err)
	}
}

func TestRunNoListedPlayers(t *testing.T) {
	r, c, _, s := createBidSubcommand()
	r.On("List").Return([]*model.Player{}, nil)
	c.On("Auth").Return(c.GetClubInfoResponse(0), nil)
	_, err := s.Run()
	if err != nil {
		t.Errorf("nil error expected, '%v' returned", err)
	}
}

type bidRunOutputTest struct {
	success []*model.Player
	failing []*model.Player
}

var (
	p1                = &model.Player{Id: 1}
	p2                = &model.Player{Id: 2}
	bidRunOutputTests = []*bidRunOutputTest{
		{[]*model.Player{}, []*model.Player{}},
		{[]*model.Player{p1}, []*model.Player{}},
		{[]*model.Player{}, []*model.Player{p1}},
		{[]*model.Player{p1}, []*model.Player{p2}},
	}
)

func TestRunPlayerBidResults(t *testing.T) {
	r, c, b, s := createBidSubcommand()
	c.On("Auth").Return(c.GetClubInfoResponse(0), nil)

	for _, tc := range bidRunOutputTests {
		all := append(tc.success, tc.failing...)
		successCount := len(tc.success)
		failingCount := len(tc.failing)
		r.On("List").Once().Return(all, nil)

		for _, p := range tc.failing {
			b.On("Bid", p, uint(0)).Once().Return(&playerbid.ErrCouldNotBid{})
		}
		for _, p := range tc.success {
			b.On("Bid", p, uint(0)).Once().Return(nil)
		}

		output, err := s.Run()
		if err != nil {
			t.Errorf("nil error expected, '%v' returned", err)
		}

		bidOutput := output.(BidSubcommandOutput)
		if bidOutput.Ok != uint(successCount) {
			t.Errorf("'%d' successful bids expected, '%d' returned", successCount, bidOutput.Ok)
		}
		if bidOutput.Failed != uint(failingCount) {
			t.Errorf("'%d' failed bids expected, '%d' returned", failingCount, bidOutput.Failed)
		}
	}
}
