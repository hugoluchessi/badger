package badger

import (
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

type middleware struct {
	fn MiddlewareFunc
}

func (m *middleware) BuildHandler(h http.Handler) http.Handler {
	return m.fn(h)
}
