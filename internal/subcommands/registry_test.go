package subcommands_test

import (
	"errors"
	"reflect"
	"sokker-org-auto-bidder/internal/subcommands"
	"testing"
)

func TestNewsSubcommandRegistryCreation(t *testing.T) {
	r := subcommands.NewSubcommandRegistry()
	if r == nil {
		t.Error("*subcommandRegistry expected but nil returned")
	}
}

func TestRegisterSubcommand(t *testing.T) {
	r := subcommands.NewSubcommandRegistry()
	r.Register("test", subcommands.NewMockSubcommand())
	sMap := reflect.ValueOf(r).Elem().FieldByName("m")
	if sMap.Kind() != reflect.Map {
		t.Fatalf("subcommand registry m field should be subcommand map")
	}

	sFound := false
	for _, s := range sMap.MapKeys() {
		if s.String() == "test" {
			sFound = true
			break
		}
	}
	if !sFound {
		t.Error("registered 'test' subcommand not found in registry")
	}	
}

// test for subcommand overwrite in registry

func TestRunNotExistingSubcommand(t *testing.T) {
	r := subcommands.NewSubcommandRegistry()
	var expectedErr *subcommands.ErrSubcommandNotAvailable
	err := r.Run("test", []string{})
	if err == nil || !errors.As(err, &expectedErr) {
		t.Errorf("expected '%T' but '%T' returned", expectedErr, err)
	}
}

// tests for run with missing flags in init
// tests for init with unrecignized error
// tests for run with run error

func TestRunExistingSubcommandWithSuccess(t *testing.T) {
	args := []string{}
	s := subcommands.NewMockSubcommand()
	s.On("Init", args).Return(nil)
	s.On("Run").Return(nil)
	r := subcommands.NewSubcommandRegistry()
	r.Register("test", s)
	err := r.Run("test", args)
	if err != nil {
		t.Errorf("expected successful run but %T returned", err)
	}
}

// tests fot getting subcommand names
