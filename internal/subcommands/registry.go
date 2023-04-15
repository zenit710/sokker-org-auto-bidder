package subcommands

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// subcommandRegistry manage available commands
type subcommandRegistry struct {
	m map[string]Subcommand
}

// NewSubcommandRegistry returns new empty registry for subcommands managment
func NewSubcommandRegistry() *subcommandRegistry {
	log.Trace("creating new subcommand registry")
	r := &subcommandRegistry{}
	r.m = make(map[string]Subcommand)
	return r
}

// Register adds subcommand to the registry on the name key
func (s *subcommandRegistry) Register(name string, cmd Subcommand) {
	log.Debugf("register new subcommand '%s'", name)
	s.m[name] = cmd
}

// Run executes subcommand registered on name key with provided args
func (s *subcommandRegistry) Run(name string, args []string) error {
	log.Tracef("trying to run '%s' subcommand", name)
	cmd := s.m[name]
	if cmd == nil {
		return &ErrSubcommandNotAvailable{Name: name, Available: s.GetSubcommandNames()}
	}
	
	log.Debugf("'%s' subcommand init with args %v", name, args)
	if err := cmd.Init(args); err != nil {
		switch err.(type) {
		case *ErrMissingFlags: return err
		default:
			log.Error(err)
			return fmt.Errorf("'%s' subcommand initialization failed", name)
		}
	}

	log.Debugf("'%s' subcommand run", name)
	return cmd.Run()
}

// GetSubcommandNames returns key names of all registered subcommands
func (s *subcommandRegistry) GetSubcommandNames() []string {
	log.Trace("get all registered subcommand names")
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
