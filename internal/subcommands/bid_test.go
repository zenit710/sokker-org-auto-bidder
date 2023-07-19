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

// type bidRunOutputTest struct {
// 	players []*model.Player
// 	failingIds []uint
// }

// var bidRunOutputTests =  []*bidRunOutputTest{
// 	{[]*model.Player{}, []uint{}},
// 	{[]*model.Player{
// 		{Id: 1},
// 	}, []uint{}},
// 	{[]*model.Player{
// 		{Id: 1},
// 	}, []uint{1}},
// 	{[]*model.Player{
// 		{Id: 1},
// 	}, []uint{2}},
// }

// func TestRunPlayerBidResults(t *testing.T) {
// 	r := player.NewMockPlayerRepository()
// 	c := client.NewMockClient()
// 	c.On("Auth").Return(c.GetEmptyClubInfoResponse(), nil)

// 	for _, tc := range bidRunOutputTests {
// 		r.On("List").Return(tc.players, nil)
// 		// mock errors from handlePlayer()
// 		// chceck every error from this method
// 		// maybe we need to extend bidRunOutputTest struct with errors for each playerID
// 		s := NewBidSubcommand(r, c)
// 		output, err := s.Run()
// 		if err != nil {
// 			t.Errorf("nil error expected, '%v' returned", err)
// 		}
// 	}
// }

// test run with player handling E2E
// test player handle errors
