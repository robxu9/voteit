package methods

import (
	"github.com/robxu9/voteit/votes"
	"github.com/satori/go.uuid"
)

// NewMajorityMethod takes 0-n votes and creates candidates 0-n with x votes.
func NewMajorityMethod(entries ...int64) *MajorityMethod {
	method := &MajorityMethod{
		Candidates: make([]uuid.UUID, len(entries)),
		Votes:      []votes.SimpleVote{},
	}
	for k, v := range entries {
		method.Candidates[k] = uuid.NewV4()
		var i int64
		for i = 0; i < v; i++ {
			method.Votes = append(method.Votes, votes.SimpleVote{
				Vote: method.Candidates[k],
			})
		}
	}
	return method
}

// MajorityMethod defines a simple majority - the candidate that achieves
// the majority of votes wins.
type MajorityMethod struct {
	Candidates []uuid.UUID
	Votes      []votes.SimpleVote

	results map[uuid.UUID]int64
	winners []uuid.UUID
}

// Options returns the candidates.
func (m *MajorityMethod) Options() []uuid.UUID {
	return m.Candidates
}

// Calculate calculates the majority vote.
func (m *MajorityMethod) Calculate() error {
	if m.results != nil {
		return ErrMethodAlreadyCalculated
	}

	m.results = map[uuid.UUID]int64{}

	var highest int64
	m.winners = []uuid.UUID{}

	for _, v := range m.Candidates {
		m.results[v] = 0
	}

	for _, v := range m.Votes {
		if _, ok := m.results[v.Vote]; !ok {
			return ErrBadVote
		}

		m.results[v.Vote]++
	}

	for k, v := range m.results {
		if v > highest {
			m.winners = []uuid.UUID{
				k,
			}
			highest = v
		} else if v == highest {
			m.winners = append(m.winners, k)
		}
	}

	if len(m.winners) == 0 {
		return ErrMethodNoWinner
	}

	if len(m.winners) > 1 {
		return ErrMethodTied
	}

	return nil
}

// Winner returns the winners.
func (m *MajorityMethod) Winner() []uuid.UUID {
	return m.winners
}
