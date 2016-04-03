package routers

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"gopkg.in/errors.v0"

	"goji.io"

	"golang.org/x/net/context"
)

var (
	ErrNotFound         = errors.New("path not found")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrBadRequest       = errors.New("bad request")
	ErrForbidden        = errors.New("forbidden")
	ErrPermission       = errors.New("permission error")
)

func PanicMiddleware(inner goji.Handler) goji.Handler {
	panicHandler := &PanicHandler{}

	return goji.HandlerFunc(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panicHandler.ServeHTTPC(err, ctx, rw, r)
			}
		}()

		inner.ServeHTTPC(ctx, rw, r)
	})
}

type PanicHandler struct{}

func (p *PanicHandler) ServeHTTPC(ex interface{}, ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	switch ex {
	case ErrNotFound:
		p.Err404(ctx, rw, r)
	case ErrMethodNotAllowed:
		p.Err405(ctx, rw, r)
	case ErrBadRequest:
		p.Err400(ctx, rw, r)
	case ErrForbidden:
		p.Err403(ctx, rw, r)
	case ErrPermission:
		p.Err550(ctx, rw, r)
	default:
		p.Err500(ex, ctx, rw, r)
	}
}

func (p *PanicHandler) Err404(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	renderer.HTML(rw, 404, "error", map[string]interface{}{
		"error": "404 Page Not Found",
	})

}

func (p *PanicHandler) Err400(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	renderer.HTML(rw, 400, "error", map[string]interface{}{
		"error": "400 Bad Request",
	})

}

func (p *PanicHandler) Err403(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	renderer.HTML(rw, 403, "error", map[string]interface{}{
		"error": "403 Unauthorised",
	})

}

func (p *PanicHandler) Err405(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	renderer.HTML(rw, 405, "error", map[string]interface{}{
		"error": "405 Method Not Allowed",
	})
}

func (p *PanicHandler) Err500(ex interface{}, ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	stackTrace := errors.Wrap(ex, 4).ErrorStack()

	logrus.WithFields(logrus.Fields{
		"err":        ex,
		"from":       r.RemoteAddr,
		"stacktrace": stackTrace,
	}).Error("Internal Server Error")

	renderer.HTML(rw, 500, "error", map[string]interface{}{
		"error": "500 Internal Server Error - see logrus",
	})
}

func (p *PanicHandler) Err550(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	renderer.HTML(rw, 550, "error", map[string]interface{}{
		"error": "550 Permissions Not Found",
	})

}
