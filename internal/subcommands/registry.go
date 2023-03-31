package subcommands

import "fmt"

type subcommandRegistry struct {
	m map[string]Subcommand
}

// NewSubcommandRegistry returns new empty registry for subcommands managment
func NewSubcommandRegistry() *subcommandRegistry {
	r := &subcommandRegistry{}
	r.m = make(map[string]Subcommand)
	return r
}

// Register adds subcommand to the registry on the name key
func (s *subcommandRegistry) Register(name string, cmd Subcommand) {
	s.m[name] = cmd
}

// Run executes subcommand registered on name key with provided args
func (s *subcommandRegistry) Run(name string, args []string) error {
	cmd := s.m[name]
	if cmd == nil {
		return &ErrSubcommandNotAvailable{Name: name, Available: s.GetSubcommandNames()}
	}
	
	if err := cmd.Init(args); err != nil {
		return err
	}

	return cmd.Run()
}

// GetSubcommandNames returns key names of all registered subcommands
func (s *subcommandRegistry) GetSubcommandNames() []string {
	names := make([]string, len(s.m))

	i := 0
	for k := range s.m {
		names[i] = k
		i++
	}

	return names
}

// ErrSubcommandNotAvailable is raised when registry has no subcommand registere on provided name key
type ErrSubcommandNotAvailable struct {
	Name string
	Available []string
}

func (e *ErrSubcommandNotAvailable) Error() string {
	return fmt.Sprintf("'%s' command is not available. Expected one of %v", e.Name, e.Available)
}
