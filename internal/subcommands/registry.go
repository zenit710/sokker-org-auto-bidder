package subcommands

import "fmt"

type subcommandRegistry struct {
	m map[string]Subcommand
}

func NewSubcommandRegistry() *subcommandRegistry {
	r := &subcommandRegistry{}
	r.m = make(map[string]Subcommand)
	return r
}

func (s *subcommandRegistry) Register(name string, cmd Subcommand) {
	s.m[name] = cmd
}

func (s *subcommandRegistry) Run(name string, args []string) error {
	cmd := s.m[name]
	if cmd == nil {
		return &ErrSubcommandNotAvailable{Name: name, Available: s.GetSubcommandNames()}
	}
	
	return cmd.Run(args)
}

func (s *subcommandRegistry) GetSubcommandNames() []string {
	names := make([]string, len(s.m))

	i := 0
	for k := range s.m {
		names[i] = k
		i++
	}

	return names
}

type ErrSubcommandNotAvailable struct {
	Name string
	Available []string
}

func (e *ErrSubcommandNotAvailable) Error() string {
	return fmt.Sprintf("'%s' command is not available. Expected one of %v", e.Name, e.Available)
}
