package playerbid

import "sokker-org-auto-bidder/internal/model"

// PlayerBidService handles player bid process
type PlayerBidService interface {
	// Bid handle player bid process
	Bid(p *model.Player, clubId uint) error
}
