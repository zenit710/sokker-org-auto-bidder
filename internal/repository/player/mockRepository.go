package player

import (
	"sokker-org-auto-bidder/internal/model"

	"github.com/stretchr/testify/mock"
)

var _ PlayerRepository = &mockPlayerRepository{}

type mockPlayerRepository struct {
	mock.Mock
}

func NewMockPlayerRepository() *mockPlayerRepository {
	return &mockPlayerRepository{}
}

func (r *mockPlayerRepository) Init() error {
	args := r.Called()
	return args.Error(0)
}

func (r *mockPlayerRepository) Add(player *model.Player) error {
	args := r.Called(player)
	return args.Error(0)
}

func (r *mockPlayerRepository) Delete(player *model.Player) error {
	args := r.Called(player)
	return args.Error(0)
}

func (r *mockPlayerRepository) List() ([]*model.Player, error) {
	args := r.Called()
	return args.Get(0).([]*model.Player), args.Error(1)
}

func (r *mockPlayerRepository) Update(player *model.Player) error {
	args := r.Called(player)
	return args.Error(0)
}
