package client

type Client interface {
	Auth() (*clubInfoResponse, error)
	Bid(id, price uint) (*transferInfoResponse, error)
	ClubInfo() (*clubInfoResponse, error)
	FetchPlayerInfo(id uint) (*playerInfoResponse, error)
}
