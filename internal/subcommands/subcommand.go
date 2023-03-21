package subcommands

type Subcommand interface {
	Run() error
}
