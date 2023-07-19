package subcommands

import "github.com/stretchr/testify/mock"

var _ Subcommand = &MockSubcommand{}

type MockSubcommand struct {
	mock.Mock
}

func (s *MockSubcommand) Init(cmdArgs []string) error {
	args := s.Called(cmdArgs)
	return args.Error(0)
}

func (s *MockSubcommand) Run() (interface{}, error) {
	args := s.Called()
	return args.Get(0), args.Error(1)
}
