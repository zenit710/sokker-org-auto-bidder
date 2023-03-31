package subcommands

type Subcommand interface {
	// Init initiaties subcommand before run with args
	Init(args []string) error
	// Run executes subcommand
	Run() error
}
