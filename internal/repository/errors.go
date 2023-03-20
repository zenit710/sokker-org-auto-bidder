package repository

type ErrRepositoryInitFailure struct {
	msg string
}

func NewErrRepositoryInitFailure(msg string) *ErrRepositoryInitFailure {
	return &ErrRepositoryInitFailure{msg}
}

func (e *ErrRepositoryInitFailure) Error() string {
	return e.msg
}