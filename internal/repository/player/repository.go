package player

import "sokker-org-auto-bidder/internal/model"

type PlayerRepository interface {
	// Init initiates repository connection
	Init() error
	// Add insert new player to the bid list
	Add(player *model.Player) error
	// Delete removes player from bid list
	Delete(player *model.Player) error
	// List returns list of all players to bid
	List() ([]*model.Player, error)
	// Update updates player in bid list
	Update(player *model.Player) error
	// Close closes repository connection
	Close()
}
