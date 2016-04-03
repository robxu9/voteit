package methods

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestMajority(t *testing.T) {
	method := NewMajorityMethod(5, 10, 15, 20, 2)

	err := method.Calculate()
	if err != nil {
		t.Fatalf("shouldn't have failed to calculate: %v", err)
	}

	winners := method.Winner()
	if len(winners) != 1 {
		t.Fatal("only one winner!")
	}

	if !uuid.Equal(winners[0], method.Candidates[3]) {
		t.Fatalf("winner is not candidate #4! (%v was winner, %v are candidates, %v are results)", winners, method.Candidates, method.results)
	}
}
