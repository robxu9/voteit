# VoteIt

VoteIt is an open-source voting application. It provides hook-ins for allowing
users to vote on selected elections (where an election could be pertaining to
anything from a simple vote to a full-scale leadership election).

## Why make this?

Voting systems are still a rarity in the open source world. There are multiple
libraries that provide various voting records, but many of them are
unmaintained. VoteIt is an attempt to implement the Schulze Method, a Condorcet
Method, into the Go language, as well as a simple majority (initially) and apply
it to different scenarios that people can vote on.

The server and client provided in this repository works with sample voting
input, such as rating Capital One ATMs within a certain area (which one is voted
the best?), local businesses within the area (such as in Foursquare - which one
is king?), or even just back to politics and electing a next leader.

## Running

Requires Go 1.6+ and its various dependencies.

Install dependencies with `go get -u -v github.com/robxu9/voteit`, and build
with `go build`. Alternatively, run directly with `go run voteit.go`.

You need a Foursquare API ID & Secret for FourSquare access. Set it in the
environment as `FOURSQUARE_ID` and `FOURSQUARE_SECRET`.

## License
[MIT licensed](http://robxu9.mit-license.org)
