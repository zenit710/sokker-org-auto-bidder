package playerbid

import (
	"sokker-org-auto-bidder/internal/model"

	"github.com/stretchr/testify/mock"
)

var _ PlayerBidService = &MockPlayerBidService{}

type MockPlayerBidService struct {
	mock.Mock
}

func (s *MockPlayerBidService) Bid(p *model.Player, clubId uint) error {
	args := s.Called(p, clubId)
	return args.Error(0)
}
