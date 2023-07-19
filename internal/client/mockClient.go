package client

import "github.com/stretchr/testify/mock"

var _ Client = &MockClient{}

type MockClient struct {
	mock.Mock
}

func (c *MockClient) GetEmptyClubInfoResponse() *clubInfoResponse {
	return &clubInfoResponse{}
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
