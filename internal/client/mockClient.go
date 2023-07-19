package client

import "github.com/stretchr/testify/mock"

var _ Client = &MockClient{}

type MockClient struct {
	mock.Mock
}

// GetClubInfoResponse returnes mocked clubInfoResponse struct
func (c *MockClient) GetClubInfoResponse(id uint) *clubInfoResponse {
	return &clubInfoResponse{Team: team{id}}
}

// GetPlayerInfoResponse returnes mocked playerInfoResponse struct
func (c *MockClient) GetPlayerInfoResponse(deadlineDate string, minBid uint, buyerId uint) *playerInfoResponse {
	return &playerInfoResponse{
		Transfer: transfer{
			Deadline: playerInfoDeadline{
				Date: deadlineDate,
			},
			Price: price{
				MinBid: bidState{
					Value: minBid,
				},
			},
			BuyerId: buyerId,
		},
	}
}

// GetTransferInfoResponse returnes mocked transferInfoResponse struct
func (c *MockClient) GetTransferInfoResponse(deadline string, minBid, buyerId uint) *transferInfoResponse {
	return &transferInfoResponse{
		Deadline: transferInfoDeadline{
			Date: date{
				Date: deadline,
			},
		},
		Price: price{
			MinBid: bidState{
				Value: minBid,
			},
		},
		Buyer: buyer{
			Id: buyerId,
		},
	}
}

func (c *MockClient) Auth() (*clubInfoResponse, error) {
	args := c.Called()
	return args.Get(0).(*clubInfoResponse), args.Error(1)
}

func (c *MockClient) Bid(id, price uint) (*transferInfoResponse, error) {
	args := c.Called(id, price)
	return args.Get(0).(*transferInfoResponse), args.Error(1)
}

func (c *MockClient) ClubInfo() (*clubInfoResponse, error) {
	args := c.Called()
	return args.Get(0).(*clubInfoResponse), args.Error(1)
}

func (c *MockClient) FetchPlayerInfo(id uint) (*playerInfoResponse, error) {
	args := c.Called(id)
	return args.Get(0).(*playerInfoResponse), args.Error(1)
}
