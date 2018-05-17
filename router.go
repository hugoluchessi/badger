package badger

import (
	"net/http"
	"sync"
)

type Router struct {
	basepath    string
	middlewares []*middleware
	routes      []Route
	lock        sync.RWMutex
}

func NewRouter(path string) *Router {
	return &Router{path, []*middleware{}, []Route{}, sync.RWMutex{}}
}

func (r *Router) Delete(path string, handler http.Handler) {
	r.Handle("DELETE", path, handler)
}

func (r *Router) Get(path string, handler http.Handler) {
	r.Handle("GET", path, handler)
}

func (r *Router) Head(path string, handler http.Handler) {
	r.Handle("HEAD", path, handler)
}

func (r *Router) Options(path string, handler http.Handler) {
	r.Handle("OPTIONS", path, handler)
}

func (r *Router) Patch(path string, handler http.Handler) {
	r.Handle("PATCH", path, handler)
}

func (r *Router) Post(path string, handler http.Handler) {
	r.Handle("POST", path, handler)
}

func (r *Router) Put(path string, handler http.Handler) {
	r.Handle("PUT", path, handler)
}

func (r *Router) Handle(method string, path string, handler http.Handler) {
	r.lock.Lock()
	defer r.lock.Unlock()

	route := Route{method, path, handler}
	r.routes = append(r.routes, route)
}

func (r *Router) Use(mwf MiddlewareFunc) {
	r.lock.Lock()
	defer r.lock.Unlock()

	mw := &middleware{mwf}

	r.middlewares = append(r.middlewares, mw)
}
