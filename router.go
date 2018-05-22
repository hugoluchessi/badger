package badger

import (
	"net/http"
	"sync"
)

// Router is responsible for gathering all routing information and to build all
// handler chaining
type Router struct {
	basepath    string
	middlewares []*middleware
	routes      []Route
	lock        sync.RWMutex
}

// NewRouter returns a pointer to a newly created router
func NewRouter(path string) *Router {
	return &Router{path, []*middleware{}, []Route{}, sync.RWMutex{}}
}

// Delete creates a new handler for DELETE method in the router
func (r *Router) Delete(path string, handler http.Handler) {
	r.Handle("DELETE", path, handler)
}

// Get creates a new handler for GET method in the router
func (r *Router) Get(path string, handler http.Handler) {
	r.Handle("GET", path, handler)
}

// Head creates a new handler for HEAD method in the router
func (r *Router) Head(path string, handler http.Handler) {
	r.Handle("HEAD", path, handler)
}

// Options creates a new handler for OPTIONS method in the router
func (r *Router) Options(path string, handler http.Handler) {
	r.Handle("OPTIONS", path, handler)
}

// Patch creates a new handler for PATCH method in the router
func (r *Router) Patch(path string, handler http.Handler) {
	r.Handle("PATCH", path, handler)
}

// Post creates a new handler for POST method in the router
func (r *Router) Post(path string, handler http.Handler) {
	r.Handle("POST", path, handler)
}

// Put creates a new handler for PUT method in the router
func (r *Router) Put(path string, handler http.Handler) {
	r.Handle("PUT", path, handler)
}

// Handle creates a new handler for the given method in the router
func (r *Router) Handle(method string, path string, handler http.Handler) {
	r.lock.Lock()
	defer r.lock.Unlock()

	route := Route{method, path, handler}
	r.routes = append(r.routes, route)
}

// Use creates a new middleware for the given functions
func (r *Router) Use(mwf MiddlewareFunc) {
	r.lock.Lock()
	defer r.lock.Unlock()

	mw := &middleware{mwf}

	r.middlewares = append(r.middlewares, mw)
}
