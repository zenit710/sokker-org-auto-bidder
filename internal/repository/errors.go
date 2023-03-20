package repository

type ErrRepositoryInitFailure struct {
	msg string
}

func (e *ErrRepositoryInitFailure) Error() string {
	return e.msg
}