package subcommands

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"testing"
)

func createSubcommand() (*client.MockClient, Subcommand) {
	c := &client.MockClient{}
	s := NewCheckAuthSubcommand(c)
	return c, s
}

func TestRunFailWhenBadCredentials(t *testing.T) {
	c, s := createSubcommand()
	c.On("Auth").Return(nil, client.ErrBadCredentials)

	if club, _ := s.Run(); club != nil {
		t.Errorf("expected club to be <nil> but '%v' returned", club)
	}
}

func TestRunFailWhenAuthRequestFailed(t *testing.T) {
	c, s := createSubcommand()
	c.On("Auth").Return(nil, errors.New(""))

	if _, err := s.Run(); err == nil || !errors.Is(err, ErrAuthFailed) {
		t.Errorf("expected '%v' error but '%v' returned", ErrAuthFailed, err)
	}
}

func TestRunFailNoClubInfo(t *testing.T) {
	c, s := createSubcommand()
	c.On("Auth").Return(nil, errors.New(""))

	if club, _ := s.Run(); club != nil {
		t.Errorf("expected club to be <nil> but '%v' returned", club)
	}
}

// func TestRunSuccess(t *testing.T) {
// 	c, s := createSubcommand()
// 	c.On("Auth").Return(c.GetClubInfoResponse(0), nil)

// 	if club, _ := s.Run(); club != nil {
// 		t.Errorf("expected club to be <nil> but '%v' returned", club)
// 	}
// }

// TOOD now we can test only that auth was failed - return better data
