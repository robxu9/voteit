package methods

import (
	"strconv"
	"strings"

	"github.com/robxu9/voteit/votes"
	"github.com/satori/go.uuid"
)

// NewSchulzeMethod provides a shorthand for initialisng a SchulzeMethod. It
// initialises new random UUIDs for the number of candidates provided, then
// uses the map to create a ranked vote for each entry provided.
//
// For example, given 4 candidates and their votes, you may call this method
// likewise:
//      method, err := NewSchulzeMethod(4, map[string]int{
//          "1,3,4,2": 3,
//          "2,1,3,4": 9,
//          "3,4,1,2": 8,
//          "4,1,2,3": 5,
//          "4,2,3,1": 5,
//      })
// Note that the string is 1-indexed but the candidates will be returned to
// you in a 0-indexed array.
func NewSchulzeMethod(candidates int, counts map[string]int) (*SchulzeMethod, error) {
	method := &SchulzeMethod{
		Candidates: make([]uuid.UUID, candidates),
		Votes:      []votes.RankedVote{},
	}

	for i := 0; i < candidates; i++ {
		method.Candidates[i] = uuid.NewV4()
	}

	for order, count := range counts {
		vote := strings.Split(order, ",")
		uuids := []uuid.UUID{}
		for _, v := range vote {
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}

			uuids = append(uuids, method.Candidates[i-1])
		}

		for i := 0; i < count; i++ {
			method.Votes = append(method.Votes, votes.RankedVote{
				Vote: uuids,
			})
		}
	}

	return method, nil
}

// SchulzeMethod represents a list of candidates and a set of votes cast
// that will be calculated.
type SchulzeMethod struct {
	Candidates []uuid.UUID
	Votes      []votes.RankedVote

	candidateIToU []uuid.UUID
	candidateUToI map[uuid.UUID]int
	provided      [][]int64
	calculated    [][]int64
	winners       map[uuid.UUID]bool
}

// Options are the valid candidates
func (s *SchulzeMethod) Options() []uuid.UUID {
	return s.Candidates
}

// Calculate calculates the strongest past from candidate to candidate via the
// Floyd-Warshall algorithm.
func (s *SchulzeMethod) Calculate() error {
	if s.calculated != nil {
		return ErrMethodAlreadyCalculated
	}

	// take a copy of it so that nobody can modify this anymore
	// also initialise the provided matrix
	s.candidateIToU = make([]uuid.UUID, len(s.Candidates))
	copy(s.candidateIToU, s.Candidates)
	s.candidateUToI = map[uuid.UUID]int{}
	s.provided = make([][]int64, len(s.Candidates))
	s.calculated = make([][]int64, len(s.Candidates))
	for k, v := range s.Candidates {
		s.candidateUToI[v] = k
		s.provided[k] = make([]int64, len(s.Candidates))
		s.calculated[k] = make([]int64, len(s.Candidates))
	}
	s.winners = map[uuid.UUID]bool{}
	for k := range s.candidateUToI {
		s.winners[k] = true
	}

	// input votes
	for _, v := range s.Votes {
		// assert that the votes are valid
		for _, vote := range v.Vote {
			if _, ok := s.candidateUToI[vote]; !ok {
				return ErrBadVote
			}
		}
		// now process it
		for i := len(v.Vote) - 2; i >= 0; i-- {
			cdd := v.Vote[i]
			for j := i + 1; j < len(v.Vote); j++ {
				s.provided[s.candidateUToI[cdd]][s.candidateUToI[v.Vote[j]]]++
			}
		}
	}

	// calculate
	s.runSchulze()

	// check for winners
	found := false
	for _, v := range s.winners {
		if v {
			if found {
				return ErrMethodTied
			}
			found = true
		}
	}
	if !found {
		return ErrMethodNoWinner
	}

	return nil
}

// Winner provides the UUIDs of the winners.
func (s *SchulzeMethod) Winner() []uuid.UUID {
	winners := []uuid.UUID{}
	for k, v := range s.winners {
		if v {
			winners = append(winners, k)
		}
	}

	return winners
}

func (s *SchulzeMethod) runSchulze() {
	size := len(s.provided)

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i != j {
				if s.provided[i][j] > s.provided[j][i] {
					s.calculated[i][j] = s.provided[i][j]
				} else {
					s.calculated[i][j] = 0
				}
			}
		}
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i != j {
				for k := 0; k < size; k++ {
					if i != k && j != k {
						s.calculated[j][k] = s.max(s.calculated[j][k], s.min(s.calculated[j][i], s.calculated[i][k]))
					}
				}
			}
		}
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i != j {
				if s.calculated[j][i] > s.calculated[i][j] {
					s.winners[s.candidateIToU[i]] = false
				}
			}
		}
	}
}

func (s *SchulzeMethod) min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func (s *SchulzeMethod) max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
