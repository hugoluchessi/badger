package badger

import (
	"net/http"
)

// Route struct defines the information needed to build a route
type Route struct {
	method  string
	path    string
	handler http.Handler
}
