package subcommands

type Subcommand interface {
	Run(args []string) error
}
