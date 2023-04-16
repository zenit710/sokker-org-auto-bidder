package session

type SessionRepository interface {
	// Get returns current session identifier
	Get() (string, error)
	// Init initalizes repository before run
	Init() error
	// Save adds new session identifier
	Save(sess string) error
}
