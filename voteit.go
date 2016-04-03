package main

import (
	"net/http"
	"os"

	"goji.io/middleware"
	"goji.io/pat"

	"github.com/Sirupsen/logrus"
	"github.com/robxu9/voteit/routers"
	"github.com/zenazn/goji/web/mutil"
	"goji.io"
	"golang.org/x/net/context"
)

func init() {
	logrus.SetOutput(os.Stderr)
}

func main() {
	mux := goji.NewMux()

	// handle middleware
	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			logrus.WithFields(logrus.Fields{
				"req":        r.RemoteAddr,
				"path":       r.RequestURI,
				"user-agent": r.UserAgent(),
				"referrer":   r.Referer(),
			}).Info("Request")
			wp := mutil.WrapWriter(rw) // proxy the rw for info later
			inner.ServeHTTPC(ctx, wp, r)
			logrus.WithFields(logrus.Fields{
				"resp":          r.RemoteAddr,
				"status":        wp.Status(),
				"bytes-written": wp.BytesWritten(),
			}).Info("Response")
		})
	})

	mux.UseC(routers.PanicMiddleware)

	mux.UseC(func(inner goji.Handler) goji.Handler {
		return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			routeFound := middleware.Pattern(ctx)

			if routeFound != nil {
				inner.ServeHTTPC(ctx, rw, r)
				return
			}

			panic(routers.ErrNotFound)
		})
	})

	// handle routes
	mux.HandleC(pat.Get("/"), goji.HandlerFunc(routers.MainRouter))
	mux.HandleC(pat.Get("/elections/"), goji.HandlerFunc(routers.ElectionsRouter))
	mux.HandleC(pat.Get("/elections/4sq"), goji.HandlerFunc(routers.FourSquareRouter))
	mux.HandleC(pat.Get("/elections/close"), goji.HandlerFunc(routers.CloseRouter))
	mux.HandleC(pat.Get("/vote/"), goji.HandlerFunc(routers.VoteRouter))
	mux.HandleC(pat.Get("/vote/push"), goji.HandlerFunc(routers.VotePushRouter))

	http.ListenAndServe(":3001", mux)
}
