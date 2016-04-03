package votes

import "github.com/satori/go.uuid"

// StarVote represents a voting scale for each candidate where each candidate
// is ranked on a scale of 1 to a configurable c.
type StarVote struct {
	Vote map[uuid.UUID]int
}
