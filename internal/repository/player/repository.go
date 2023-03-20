package player

import "sokker-org-auto-bidder/internal/model"

type PlayerRepository interface {
	Init() error
	Add(player *model.Player) error
	GetList() ([]*model.Player, error)
	Close()
}
