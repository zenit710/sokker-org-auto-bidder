package subcommands_test

import (
	"errors"
	"reflect"
	"sokker-org-auto-bidder/internal/subcommands"
	"testing"
	"unsafe"
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

func TestRegisterSubcommandOverwrite(t *testing.T) {
	s1 := subcommands.NewMockSubcommand()
	s2 := subcommands.NewMockSubcommand()
	s2Addr := uintptr(unsafe.Pointer(s2))
	r := subcommands.NewSubcommandRegistry()
	r.Register("test", s1)
	r.Register("test", s2)
	sMap := reflect.ValueOf(r).Elem().FieldByName("m")
	if sMap.Kind() != reflect.Map {
		t.Fatalf("subcommand registry m field should be subcommand map")
	}

	for _, s := range sMap.MapKeys() {
		if s.String() == "test" {
			sAddr := sMap.MapIndex(s).Elem().Pointer()
			if sAddr != s2Addr {
				t.Errorf("expected subcommand address same as second registered cmd (%v), but found %v", s2Addr, sAddr)
			}
			break
		}
	}
}

func TestRunNotExistingSubcommand(t *testing.T) {
	r := subcommands.NewSubcommandRegistry()
	var expectedErr *subcommands.ErrSubcommandNotAvailable
	err := r.Run("test", []string{})
	if err == nil || !errors.As(err, &expectedErr) {
		t.Errorf("expected '%T' but '%T' returned", expectedErr, err)
	}
}

func TestRunMissingExecFlags(t *testing.T) {
	expectedErr := &subcommands.ErrMissingFlags{}
	args := []string{}
	s := subcommands.NewMockSubcommand()
	s.On("Init", args).Return(expectedErr)
	r := subcommands.NewSubcommandRegistry()
	r.Register("test", s)
	
	err := r.Run("test", args)
	if err == nil || !errors.As(err, &expectedErr) {
		t.Errorf("expected '%T' but '%T' returned", expectedErr, err)
	}
}

func TestRunSubcommandInitFailed(t *testing.T) {
	args := []string{}
	s := subcommands.NewMockSubcommand()
	s.On("Init", args).Return(errors.New("error"))
	r := subcommands.NewSubcommandRegistry()
	r.Register("test", s)
	
	err := r.Run("test", args)
	if err == nil || !errors.Is(err, subcommands.ErrSubcommandInitFailed) {
		t.Errorf("expected '%v' but '%v' returned", subcommands.ErrSubcommandInitFailed, err)
	}
}

func TestRunSubcommandRunFailed(t *testing.T) {
	expectedErr := errors.New("error")
	args := []string{}
	s := subcommands.NewMockSubcommand()
	s.On("Init", args).Return(nil)
	s.On("Run").Return(expectedErr)
	r := subcommands.NewSubcommandRegistry()
	r.Register("test", s)
	
	err := r.Run("test", args)
	if err == nil || !errors.Is(err, expectedErr) {
		t.Errorf("expected '%v' but '%v' returned", expectedErr, err)
	}
}

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

type subcommandNamesTest struct {
	names []string
	expectedCount int
}

var subcommandNamesTests = []subcommandNamesTest{
	{[]string{}, 0},
	{[]string{"test"}, 1},
	{[]string{"test", "test2"}, 2},
	{[]string{"test", "test2", "test"}, 2},
}

func TestGetSubcommandNames(t *testing.T) {
	for i, test := range subcommandNamesTests {
		r := subcommands.NewSubcommandRegistry()
		for _, name := range test.names {
			r.Register(name, subcommands.NewMockSubcommand())
		}

		rNames := r.GetSubcommandNames()
		rNamesLen := len(rNames)
		if rNamesLen != test.expectedCount {
			t.Errorf("expected %v subcommand registers but found %v in test set index %v", test.expectedCount, rNamesLen, i)
		}
	}
}
