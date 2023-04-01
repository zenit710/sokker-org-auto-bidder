package repository

/*
ErrRepositoryInitFailure is raised when something goes wrong with repository Init.
It can be DB connection error for instance.
*/
type ErrRepositoryInitFailure struct {
	msg string
}

// NewErrRepositoryInitFailure creates new repository Init error with message
func NewErrRepositoryInitFailure(msg string) *ErrRepositoryInitFailure {
	return &ErrRepositoryInitFailure{msg}
}

func (e *ErrRepositoryInitFailure) Error() string {
	return e.msg
}