package client

type Client interface {
	// Auth authorize user in game
	Auth() (*clubInfoResponse, error)
	// Bid makes player bid
	Bid(id, price uint) (*transferInfoResponse, error)
	// ClubInfo returns info about user club
	ClubInfo() (*clubInfoResponse, error)
	// FetchPlayerInfo returns info about player
	FetchPlayerInfo(id uint) (*playerInfoResponse, error)
}
