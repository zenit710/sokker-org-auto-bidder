package subcommands

import "github.com/stretchr/testify/mock"

var _ Subcommand = &mockSubcommand{}

type mockSubcommand struct {
	mock.Mock
}

func NewMockSubcommand() *mockSubcommand {
	return &mockSubcommand{}
}

func (s *mockSubcommand) Init(cmdArgs []string) error {
	args := s.Called(cmdArgs)
	return args.Error(0)
}

func (s *mockSubcommand) Run() (interface{}, error) {
	args := s.Called()
	return args.Get(0), args.Error(1)
}
