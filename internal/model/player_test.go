package model_test

import (
	"sokker-org-auto-bidder/internal/model"
	"testing"
)

type playerTest struct {
	player *model.Player
	valid  bool
}

func TestValidate(t *testing.T) {
	tests := []playerTest{
		{&model.Player{Id: 0, MaxPrice: 0}, false},
		{&model.Player{Id: 0, MaxPrice: 1}, false},
		{&model.Player{Id: 1, MaxPrice: 0}, false},
		{&model.Player{Id: 1, MaxPrice: 1}, true},
	}

	for i, test := range tests {
		if err := test.player.Validate(); (err == nil) != test.valid {
			text := "valid"
			if !test.valid {
				text = "invalid"
			}
			t.Errorf("expected player (test case no. %d) to be %s but '%v' error returned", i, text, err)
		}
	}
}
