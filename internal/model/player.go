package model

import "time"

// Player is a struct for player in bid list
type Player struct {
	Id       uint
	MaxPrice uint
	Deadline time.Time
}

// Validate validates player struct values
func (p *Player) Validate() error {
	if p.Id <= 0 {
		return &ErrInvalidPlayerParam{"playerId has to be greater than zero"}
	}

	if p.MaxPrice <= 0 {
		return &ErrInvalidPlayerParam{"maxPrice has to be greater than zero"}
	}

	return nil
}

// ErrInvalidPlayerParam is raised when validatio error occurs
type ErrInvalidPlayerParam struct {
	msg string
}

func (e *ErrInvalidPlayerParam) Error() string {
	return e.msg
}
