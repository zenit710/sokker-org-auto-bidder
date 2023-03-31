package player

import "sokker-org-auto-bidder/internal/model"

type PlayerRepository interface {
	Init() error
	Add(player *model.Player) error
	List() ([]*model.Player, error)
	Update(player *model.Player) error
	Close()
}
