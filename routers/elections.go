package routers

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"menteslibres.net/gosexy/dig"

	"golang.org/x/net/context"

	fsq "github.com/elbuo8/4square-venues"
	"github.com/robxu9/voteit/methods"
	"github.com/robxu9/voteit/votes"
	"github.com/satori/go.uuid"
)

var (
	Places         = make(map[uuid.UUID]string)
	Elections      = make(map[uuid.UUID]*methods.SchulzeMethod)
	ElectionsMutex = &sync.RWMutex{}
)

// MainRouter handles the main page for VoteIt.
func MainRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	renderer.HTML(w, 200, "index", map[string]interface{}{})
}

// ElectionsRouter allows creation of a couple of sample elections.
func ElectionsRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ElectionsMutex.RLock()
	defer ElectionsMutex.RUnlock()

	renderer.HTML(w, 200, "elections/index", map[string]interface{}{})
}

// 4SQRouter searches FourSquare for places specified and takes the top five for voting.
func FourSquareRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fs := fsq.NewFSVenuesClient(os.Getenv("FOURSQUARE_ID"), os.Getenv("FOURSQUARE_SECRET"))
	params := map[string]string{
		"near":  r.FormValue("loc"),
		"query": r.FormValue("for"),
		"limit": "5",
	}

	venues, err := fs.GetVenues(params)
	if err != nil {
		panic(err)
	}

	castVenue := venues.(map[string]interface{})

	resultCode := dig.Int64(&castVenue, "meta", "code")
	if resultCode != 200 {
		panic(ErrBadRequest)
	}

	// open up an election
	ElectionsMutex.Lock()
	defer ElectionsMutex.Unlock()

	electionUUID := uuid.NewV4()

	candidates := []uuid.UUID{}

	for _, v := range dig.Interface(&castVenue, "response", "venues").([]interface{}) {
		castV := v.(map[string]interface{})
		name := fmt.Sprintf("%s (%s)", dig.String(&castV, "name"), dig.String(&castV, "location", "crossStreet"))
		locUUID := uuid.NewV4()
		candidates = append(candidates, locUUID)
		Places[locUUID] = name
	}

	Elections[electionUUID] = &methods.SchulzeMethod{
		Candidates: candidates,
		Votes:      []votes.RankedVote{},
	}

	dataRenderer := map[string]interface{}{
		"venues": dig.Interface(&castVenue, "response", "venues"),
		"uuid":   electionUUID.String(),
	}

	renderer.HTML(w, 200, "elections/4sq", dataRenderer)
}

// CloseRouter closes an election, calculates it, and goes ahead and displays
// the results.
func CloseRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get the uuid
	targetUUID := r.FormValue("uuid")
	toUUID := uuid.FromStringOrNil(targetUUID)

	// close
	ElectionsMutex.Lock()
	election := Elections[toUUID]
	delete(Elections, toUUID)
	ElectionsMutex.Unlock()

	if election == nil {
		panic(ErrBadRequest)
	}

	election.Calculate()

	// get names
	ElectionsMutex.RLock()
	var names []string
	for _, v := range election.Winner() {
		names = append(names, Places[v])
	}
	ElectionsMutex.RUnlock()

	renderer.HTML(w, 200, "elections/closed", names)
}

// VoteRouter allows a user to vote on an election.
func VoteRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get uuid if any
	targetUUID := r.FormValue("uuid")
	toUUID := uuid.FromStringOrNil(targetUUID)

	ElectionsMutex.RLock()
	defer ElectionsMutex.RUnlock()
	election := Elections[toUUID]

	if election == nil {
		renderer.HTML(w, 200, "vote/index", nil)
		return
	}

	// get names
	nameToUUID := map[string]string{}
	for _, v := range election.Candidates {
		nameToUUID[Places[v]] = v.String()
	}

	renderer.HTML(w, 200, "vote/vote", map[string]interface{}{
		"uuid":    targetUUID,
		"options": nameToUUID,
	})
}

// VotePushRouter allows a user to push a vote.
func VotePushRouter(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	targetUUID := r.FormValue("uuid")
	toUUID := uuid.FromStringOrNil(targetUUID)

	ElectionsMutex.Lock()
	defer ElectionsMutex.Unlock()
	election := Elections[toUUID]

	if election == nil {
		panic(ErrBadRequest)
	}

	// get a vote
	rankedVote := votes.RankedVote{
		Vote: []uuid.UUID{},
	}

	order := map[string]uuid.UUID{}

	for k, v := range r.Form {
		if k == "uuid" {
			continue
		}

		optionUUID := uuid.FromStringOrNil(k)
		order[v[0]] = optionUUID
	}

	rankedVote.Vote = append(rankedVote.Vote, order["1"], order["2"], order["3"], order["4"], order["5"])

	// add it in
	election.Votes = append(election.Votes, rankedVote)

	renderer.HTML(w, 200, "vote/voted", nil)
}
