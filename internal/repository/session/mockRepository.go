package session

import "github.com/stretchr/testify/mock"

var _ SessionRepository = &mockSessionRepository{}

type mockSessionRepository struct {
	mock.Mock
}

func NewMockSessionRepository() *mockSessionRepository {
	return &mockSessionRepository{}
}

func (r *mockSessionRepository) Get() (string, error) {
	args := r.Called()
	return args.String(0), args.Error(1)
}

func (r *mockSessionRepository) Init() error {
	args := r.Called()
	return args.Error(0)
}

func (r *mockSessionRepository) Save(sess string) error {
	args := r.Called(sess)
	return args.Error(0)
}
