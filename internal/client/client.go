package client

type Client interface {
	Auth() error
	FetchPlayerInfo(id uint) (*playerInfoResponse, error)
}
