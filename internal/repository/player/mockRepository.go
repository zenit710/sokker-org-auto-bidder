package player

import (
	"sokker-org-auto-bidder/internal/model"

	"github.com/stretchr/testify/mock"
)

var _ PlayerRepository = &MockPlayerRepository{}

type MockPlayerRepository struct {
	mock.Mock
}

func (r *MockPlayerRepository) Init() error {
	args := r.Called()
	return args.Error(0)
}

func (r *MockPlayerRepository) Add(player *model.Player) error {
	args := r.Called(player)
	return args.Error(0)
}

func (r *MockPlayerRepository) Delete(player *model.Player) error {
	args := r.Called(player)
	return args.Error(0)
}

func (r *MockPlayerRepository) List() ([]*model.Player, error) {
	args := r.Called()
	return args.Get(0).([]*model.Player), args.Error(1)
}

func (r *MockPlayerRepository) Update(player *model.Player) error {
	args := r.Called(player)
	return args.Error(0)
}
