package methods

import (
	"errors"

	"github.com/satori/go.uuid"
)

var (
	// ErrMethodAlreadyCalculated indicates that the method has been calculated already.
	ErrMethodAlreadyCalculated = errors.New("methods: method already calculated")

	// ErrMethodTied indicates that there was a tie.
	ErrMethodTied = errors.New("methods: method tied; calculate manually")

	// ErrMethodNoWinner indicates that there was no winner.
	ErrMethodNoWinner = errors.New("methods: method could not find a winner")

	// ErrBadVote indicates that there was a bad vote and processing has stopped.
	ErrBadVote = errors.New("methods: bad vote encountered, stopped processing")
)

// Method is a method of calculating the winner of votes.
type Method interface {
	// Calculate calls the method to calculate
	Calculate() error
	// Options are the valid candidates
	Options() []uuid.UUID
	// Winner reveals the winner. Only available after Calculate().
	Winner() []uuid.UUID
}
