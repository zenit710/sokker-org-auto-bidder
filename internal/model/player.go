package model

type Player struct {
	Id uint
	MaxPrice uint
	Deadline string
}

func (p *Player) Validate() error {
	if p.Id <= 0 {
		return &ErrInvalidPlayerParam{"playerId has to be greater than zero"}
	}

	if p.MaxPrice <= 0 {
		return &ErrInvalidPlayerParam{"maxPrice has to be greater than zero"}
	}

	return nil
}

type ErrInvalidPlayerParam struct {
	msg string
}

func (e *ErrInvalidPlayerParam) Error() string {
	return e.msg
}
