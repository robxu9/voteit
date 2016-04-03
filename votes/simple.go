package votes

import "github.com/satori/go.uuid"

// SimpleVote is single-candidate vote.
type SimpleVote struct {
	Vote uuid.UUID
}
