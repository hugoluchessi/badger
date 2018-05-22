package badger

import (
	"net/http"
)

// MiddlewareFunc is the interface needed to send as you Use middlewares
// in your router.
type MiddlewareFunc func(http.Handler) http.Handler

type middleware struct {
	fn MiddlewareFunc
}

func (m *middleware) BuildHandler(h http.Handler) http.Handler {
	return m.fn(h)
}
