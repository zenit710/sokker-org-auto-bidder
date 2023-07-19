package playerbid

import "fmt"

type ErrTransferInfoFetchFailed struct {
	id uint
}

func (e *ErrTransferInfoFetchFailed) Error() string {
	return fmt.Sprintf("could not fetch player (%d) transfer info", e.id)
}

type ErrMaxPriceReached struct {
	id uint
}

func (e *ErrMaxPriceReached) Error() string {
	return fmt.Sprintf("max price reached, cannot bid player (%d) further", e.id)
}

type ErrCurrentLeader struct {
	id uint
}

func (e *ErrCurrentLeader) Error() string {
	return fmt.Sprintf("you are current leader, no reason to bid player (%d)", e.id)
}

type ErrCouldNotBid struct {
	id uint
	e  error
}

func (e *ErrCouldNotBid) Error() string {
	return fmt.Sprintf("player (%d) bid could not be made: %v", e.id, e.e)
}
