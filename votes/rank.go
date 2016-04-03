package votes

import "github.com/satori/go.uuid"

// RankedVote represents a vote based on a ranking system, with most preferred
// at index zero and least preferred at index n, where there are n candidates.
type RankedVote struct {
	Vote []uuid.UUID
}
