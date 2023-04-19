package client

import "github.com/stretchr/testify/mock"

var _ Client = &mockClient{}

type mockClient struct {
	mock.Mock
}

func NewMockClient() *mockClient {
	return &mockClient{}
}

func (c *mockClient) Auth() (*clubInfoResponse, error) {
	args := c.Called()
	return args.Get(0).(*clubInfoResponse), args.Error(1)
}

func (c *mockClient) Bid(id, price uint) (*transferInfoResponse, error) {
	args := c.Called(id, price)
	return args.Get(0).(*transferInfoResponse), args.Error(1)
}

func (c *mockClient) ClubInfo() (*clubInfoResponse, error) {
	args := c.Called()
	return args.Get(0).(*clubInfoResponse), args.Error(1)
}

func (c *mockClient) FetchPlayerInfo(id uint) (*playerInfoResponse, error) {
	args := c.Called(id)
	return args.Get(0).(*playerInfoResponse), args.Error(1)
}
