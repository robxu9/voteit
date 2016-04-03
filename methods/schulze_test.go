package methods

import (
	"testing"

	"github.com/robxu9/voteit/votes"
	"github.com/satori/go.uuid"
)

func TestSchulzeNoWinner(t *testing.T) {
	method := &SchulzeMethod{
		Candidates: []uuid.UUID{},
		Votes:      []votes.RankedVote{},
	}

	err := method.Calculate()
	if err != ErrMethodNoWinner {
		t.Fatal(err)
	}

	if len(method.Winner()) != 0 {
		t.Fatal("winners isn't zero!")
	}
}

func TestSchulzeWikipedia1(t *testing.T) {

	method, err := NewSchulzeMethod(4, map[string]int{
		"1,3,4,2": 3,
		"2,1,3,4": 9,
		"3,4,1,2": 8,
		"4,1,2,3": 5,
		"4,2,3,1": 5,
	})

	if err != nil {
		t.Fatalf("failed to parse new schulze method: %v", err)
	}

	err = method.Calculate()
	if err != nil {
		t.Fatalf("should have no error: %v", err)
	}

	winner := method.Winner()
	if len(winner) != 1 {
		t.Fatalf("winner is not len == 1: %v", winner)
	}

	if winner[0] != method.Candidates[2] {
		t.Fatal("winner was not candidate #3")
	}
}

func TestSchulzeWikipedia2(t *testing.T) {

	method, err := NewSchulzeMethod(5, map[string]int{
		"1,3,2,5,4": 5,
		"1,4,5,3,2": 5,
		"2,5,4,1,3": 8,
		"3,1,2,5,4": 3,
		"3,1,5,2,4": 7,
		"3,2,1,4,5": 2,
		"4,3,5,2,1": 7,
		"5,2,1,4,3": 8,
	})

	if err != nil {
		t.Fatalf("failed to parse new schulze method: %v", err)
	}

	err = method.Calculate()
	if err != nil {
		t.Fatalf("should have no error: %v", err)
	}

	winner := method.Winner()
	if len(winner) != 1 {
		t.Fatalf("winner is not len == 1: %v", winner)
	}

	if winner[0] != method.Candidates[4] {
		t.Fatal("winner was not candidate #5")
	}

}

func TestSchulzeTie(t *testing.T) {
	method, err := NewSchulzeMethod(3, map[string]int{
		"1,2,3": 1,
		"2,1,3": 1,
	})

	if err != nil {
		t.Fatalf("failed to parse new schulze method: %v", err)
	}

	err = method.Calculate()
	if err != ErrMethodTied {
		t.Fatalf("err should be tied: %v", err)
	}

	winner := method.Winner()
	if len(winner) != 2 {
		t.Fatalf("winner is not len == 2: %v", winner)
	}

	if (winner[0] != method.Candidates[0] || winner[1] != method.Candidates[1]) && (winner[0] != method.Candidates[1] || winner[1] != method.Candidates[0]) {
		t.Fatal("winners were not #1 and #2")
	}
}
