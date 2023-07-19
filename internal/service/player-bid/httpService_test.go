package playerbid_test

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	playerbid "sokker-org-auto-bidder/internal/service/player-bid"
	"testing"
)

func createService() (*player.MockPlayerRepository, *client.MockClient, playerbid.PlayerBidService) {
	r := &player.MockPlayerRepository{}
	c := &client.MockClient{}
	s := playerbid.NewHttpPlayerBidService(r, c)
	return r, c, s	
}

func TestBidFailedWhenCanNotFetchPlayerInfo(t *testing.T) {
	p := &model.Player{Id: 1}
	var expectedErrType *playerbid.ErrTransferInfoFetchFailed

	_, c, s := createService()
	c.On("FetchPlayerInfo", p.Id).Return(c.GetEmptyPlayerInfoResponse(), errors.New(""))

	if err := s.Bid(p, 0); err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

// TODO change mockClient GetEmptyClubInfoResponse and GetEmptyPlayerInfoResponse
// allow to create objects based on method arguments
