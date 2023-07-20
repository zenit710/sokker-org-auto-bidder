package subcommands

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	"testing"

	"github.com/stretchr/testify/mock"
)

func createPlayerAddSubcommand() (*player.MockPlayerRepository, *client.MockClient, Subcommand) {
	r := &player.MockPlayerRepository{}
	c := &client.MockClient{}
	s := NewPlayerAddSubcommand(r, c)
	return r, c, s
}

func TestInitFailureWhenCanNotParseFlags(t *testing.T) {
	_, _, s := createPlayerAddSubcommand()
	if err := s.Init([]string{"--="}); err == nil || !errors.Is(err, ErrCanNotParseArgs) {
		t.Errorf("expected '%v' error but '%v' returned", ErrCanNotParseArgs, err)
	}
}

func TestInitFailureWhenMissingFlags(t *testing.T) {
	var expectedError *ErrMissingFlags
	_, _, s := createPlayerAddSubcommand()
	if err := s.Init([]string{}); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestInitSuccess(t *testing.T) {
	_, _, s := createPlayerAddSubcommand()
	if err := s.Init([]string{"--playerId", "1", "--maxPrice", "1"}); err != nil {
		t.Errorf("expected <nil> but '%v' returned", err)
	}
}

func TestRunFailureWhenCouldNotFetchTransferDetails(t *testing.T) {
	var expectedError *ErrCanNotFetchTransferDetails
	_, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "1", "--maxPrice", "1"})
	c.On("FetchPlayerInfo", uint(1)).Return(nil, errors.New(""))
	if _, err := s.Run(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestRunFailureWhenMaxPriceIsToLowToBid(t *testing.T) {
	var expectedError *ErrMaxPriceToLow
	_, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "1", "--maxPrice", "1"})
	c.On("FetchPlayerInfo", uint(1)).Return(c.GetPlayerInfoResponse("", 2, 1), nil)
	if _, err := s.Run(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestRunFailureWhenCanNotParseTransferDeadline(t *testing.T) {
	var expectedError *ErrCanNotParseTransferDeadline
	_, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "1", "--maxPrice", "1"})
	c.On("FetchPlayerInfo", uint(1)).Return(c.GetPlayerInfoResponse("", 1, 1), nil)
	if _, err := s.Run(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestRunFailureWhenInvalidPlayerId(t *testing.T) {
	var expectedError *model.ErrInvalidPlayerParam
	_, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "0", "--maxPrice", "1"})
	c.On("FetchPlayerInfo", uint(0)).Return(c.GetPlayerInfoResponse(client.TimeLayout, 1, 1), nil)
	if _, err := s.Run(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestRunFailureWhenInvalidMaxPrice(t *testing.T) {
	var expectedError *model.ErrInvalidPlayerParam
	_, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "1", "--maxPrice", "0"})
	c.On("FetchPlayerInfo", uint(1)).Return(c.GetPlayerInfoResponse(client.TimeLayout, 0, 1), nil)
	if _, err := s.Run(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestRunFailureWhenCanNotAddToBidList(t *testing.T) {
	var expectedError *ErrCanNotAddToBidList
	r, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "1", "--maxPrice", "1"})
	c.On("FetchPlayerInfo", uint(1)).Return(c.GetPlayerInfoResponse(client.TimeLayout, 0, 1), nil)
	r.On("Add", mock.Anything).Return(errors.New(""))
	if _, err := s.Run(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%v' error but '%v' returned", expectedError, err)
	}
}

func TestRunSuccess(t *testing.T) {
	r, c, s := createPlayerAddSubcommand()
	s.Init([]string{"--playerId", "1", "--maxPrice", "1"})
	c.On("FetchPlayerInfo", uint(1)).Return(c.GetPlayerInfoResponse(client.TimeLayout, 0, 1), nil)
	r.On("Add", mock.Anything).Return(nil)
	if _, err := s.Run(); err != nil {
		t.Errorf("expected <nil> but '%v' returned", err)
	}
}
