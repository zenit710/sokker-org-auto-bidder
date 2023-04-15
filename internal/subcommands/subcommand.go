package subcommands

import "fmt"

type Subcommand interface {
	// Init initiaties subcommand before run with args
	Init(args []string) error
	// Run executes subcommand
	Run() error
}

type ErrMissingFlags struct {
	Required []string
}

func (e *ErrMissingFlags) Error() string {
	return fmt.Sprintf("one of required subcommand flags missing: %v", e.Required)
}
