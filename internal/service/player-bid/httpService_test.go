package playerbid_test

import (
	"errors"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/model"
	"sokker-org-auto-bidder/internal/repository/player"
	playerbid "sokker-org-auto-bidder/internal/service/player-bid"
	"testing"
	"time"
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
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, 1, 1),
		errors.New(""),
	)

	if err := s.Bid(p, 0); err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

func TestBidFailedWhenMaxPriceReached(t *testing.T) {
	p := &model.Player{Id: 1, MaxPrice: 1}
	var expectedErrType *playerbid.ErrMaxPriceReached

	r, c, s := createService()
	r.On("Delete", p).Return(nil)
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, 2, 1),
		nil,
	)

	if err := s.Bid(p, 0); err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

func TestBidFailedWhenCurrentLeader(t *testing.T) {
	p := &model.Player{Id: 1, MaxPrice: 1}
	var userClubId uint = 1
	var expectedErrType *playerbid.ErrCurrentLeader

	_, c, s := createService()
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, 1, userClubId),
		nil,
	)

	if err := s.Bid(p, userClubId); err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

func TestBidFailedWhenBidRequestFailed(t *testing.T) {
	p := &model.Player{Id: 1, MaxPrice: 1}
	var minBid uint = 1
	var expectedErrType *playerbid.ErrCouldNotBid

	_, c, s := createService()
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, minBid, 1),
		nil,
	)
	c.On("Bid", p.Id, minBid).Return(
		c.GetTransferInfoResponse(client.TimeLayout, 1, 1),
		errors.New(""),
	)

	if err := s.Bid(p, 0); err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected '%T' error, got '%T'", expectedErrType, err)
	}
}

func TestBidPlayerDeleteFromListWhenPriceToHigh(t *testing.T) {
	p := &model.Player{Id: 1, MaxPrice: 1}

	r, c, s := createService()
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, 2, 1),
		nil,
	)
	r.On("Delete", p).Return(nil)

	s.Bid(p, 0)
	r.AssertCalled(t, "Delete", p)
}

func TestBidSuccess(t *testing.T) {
	deadline, _ := time.Parse(client.TimeLayout, client.TimeLayout)
	p := &model.Player{Id: 1, MaxPrice: 1, Deadline: deadline}
	var minBid uint = 1

	_, c, s := createService()
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, minBid, 1),
		nil,
	)
	c.On("Bid", p.Id, minBid).Return(
		c.GetTransferInfoResponse(client.TimeLayout, 1, 1),
		nil,
	)

	if err := s.Bid(p, 0); err != nil {
		t.Errorf("no errors should be raised but got '%v'", err)
	}
}

func TestBidSuccessUpdateTransferDeadlineWhenChanged(t *testing.T) {
	p := &model.Player{Id: 1, MaxPrice: 1}
	var minBid uint = 1

	r, c, s := createService()
	c.On("FetchPlayerInfo", p.Id).Return(
		c.GetPlayerInfoResponse(client.TimeLayout, minBid, 1),
		nil,
	)
	c.On("Bid", p.Id, minBid).Return(
		c.GetTransferInfoResponse(client.TimeLayout, 1, 1),
		nil,
	)
	r.On("Update", p).Return(nil)

	s.Bid(p, 0)
	r.AssertCalled(t, "Update", p)
}
