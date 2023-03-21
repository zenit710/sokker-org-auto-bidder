package subcommands

type Subcommand interface {
	Init(args []string) error
	Run() error
}
