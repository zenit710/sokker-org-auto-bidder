package subcommands

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	playerbid "sokker-org-auto-bidder/internal/service/player-bid"
	"testing"

	"github.com/stretchr/testify/mock"
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
	c.On("Auth").Return(c.GetEmptyClubInfoResponse(), errors.New("error"))
	_, err := s.Run()
	if err == nil || !errors.Is(err, ErrApiAuthFailed) {
		t.Errorf("'%v' expected, '%v' returned", ErrApiAuthFailed, err)
	}
}

func TestRunNoListedPlayers(t *testing.T) {
	r, c, _, s := createBidSubcommand()
	r.On("List").Return([]*model.Player{}, nil)
	c.On("Auth").Return(c.GetEmptyClubInfoResponse(), nil)
	_, err := s.Run()
	if err != nil {
		t.Errorf("nil error expected, '%v' returned", err)
	}
}

type bidRunOutputTest struct {
	players []*model.Player
	failingPlayers []*model.Player
	success uint
	failed uint
}

var (
	p1 = &model.Player{Id: 1}
	p2 = &model.Player{Id: 2}
	bidRunOutputTests =  []*bidRunOutputTest{
		{[]*model.Player{}, []*model.Player{}, 0, 0},
		{[]*model.Player{p1}, []*model.Player{}, 1, 0},
		{[]*model.Player{p1}, []*model.Player{p1}, 0, 1},
		{[]*model.Player{p1, p2}, []*model.Player{p1}, 1, 1},
	}
)

func TestRunPlayerBidResults(t *testing.T) {
	r, c, b, s := createBidSubcommand()
	c.On("Auth").Return(c.GetEmptyClubInfoResponse(), nil)
	// TODO mocking anything causes failing tests
	// find a way to mock every player
	// maybe create arrays not for all and falling but succesfull and falling separately
	b.On("Bid", mock.Anything, mock.Anything).Return(nil)

	for _, tc := range bidRunOutputTests {
		r.On("List").Once().Return(tc.players, nil)

		for _, p := range tc.failingPlayers {
			b.On("Bid", p, 0).Once().Return()
		}
		
		output, err := s.Run()
		if err != nil {
			t.Errorf("nil error expected, '%v' returned", err)
		}
		
		bidOutput := output.(BidSubcommandOutput)
		if bidOutput.Ok != tc.success {
			t.Errorf("'%d' successful bids expected, '%d' returned", tc.success, bidOutput.Ok)
		}
		if bidOutput.Failed != tc.failed {
			t.Errorf("'%d' failed bids expected, '%d' returned", tc.failed, bidOutput.Failed)
		}
	}
}

// test run with player handling E2E
// test player handle errors
