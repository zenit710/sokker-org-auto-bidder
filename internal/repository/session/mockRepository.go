package session

import "github.com/stretchr/testify/mock"

var _ SessionRepository = &MockSessionRepository{}

type MockSessionRepository struct {
	mock.Mock
}

func (r *MockSessionRepository) Get() (string, error) {
	args := r.Called()
	return args.String(0), args.Error(1)
}

func (r *MockSessionRepository) Init() error {
	args := r.Called()
	return args.Error(0)
}

func (r *MockSessionRepository) Save(sess string) error {
	args := r.Called(sess)
	return args.Error(0)
}
