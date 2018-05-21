package badger

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type HandleFunc func(http.Handler) httprouter.Handle

type Mux struct {
	routers          []*Router
	mainrouter       *httprouter.Router
	lock             sync.RWMutex
	NotFound         http.HandlerFunc
	MethodNotAllowed http.HandlerFunc
	PanicHandler     func(http.ResponseWriter, *http.Request, interface{})
}

func NewMux() *Mux {
	return &Mux{[]*Router{}, nil, sync.RWMutex{}, nil, nil, nil}
}

func (mux *Mux) AddRouter(path string) *Router {
	mux.lock.Lock()
	defer mux.lock.Unlock()

	router := NewRouter(path)
	mux.routers = append(mux.routers, router)

	return router
}

func (mux *Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	req.URL.Path = normalizeRoutePath(req.URL.Path)
	mux.getMainRouterInstance().ServeHTTP(res, req)
}

func (mux *Mux) getMainRouterInstance() *httprouter.Router {
	if mux.mainrouter == nil {
		mux.createMainRouterInstance()
	}

	return mux.mainrouter
}

func (mux *Mux) createMainRouterInstance() {
	mux.lock.Lock()
	defer mux.lock.Unlock()

	mux.mainrouter = httprouter.New()

	if mux.NotFound != nil {
		mux.mainrouter.NotFound = mux.NotFound
	}

	if mux.MethodNotAllowed != nil {
		mux.mainrouter.MethodNotAllowed = mux.MethodNotAllowed
	}

	if mux.PanicHandler != nil {
		mux.mainrouter.PanicHandler = mux.PanicHandler
	}

	for _, router := range mux.routers {
		routerroutes := router.buildRoutes()

		for _, route := range routerroutes {
			mux.mainrouter.Handle(
				route.method,
				route.path,
				// FIXME: Ignore params for now
				HandleFunc(func(h http.Handler) httprouter.Handle {
					return func(res http.ResponseWriter, req *http.Request, rps httprouter.Params) {
						typed := CreateRouteParams(rps)
						ctx := req.Context()
						ctx = context.WithValue(ctx, RouteParamsKey, typed)
						req = req.WithContext(ctx)

						h.ServeHTTP(res, req)
					}
				})(route.handler),
			)
		}
	}
}

func (r *Router) buildRoutes() []Route {
	builtroutes := make([]Route, 0)

	for _, route := range r.routes {
		// Ensure path starts and ends with /
		p := normalizeRoutePath(r.basepath, route.path)
		builtroute := Route{route.method, p, route.handler}

		for _, middleware := range r.middlewares {
			builtroute.handler = middleware.BuildHandler(builtroute.handler)
		}

		builtroutes = append(builtroutes, builtroute)
	}

	return builtroutes
}

func normalizeRoutePath(p ...string) string {
	rp := strings.Join(p, "/")
	return fmt.Sprintf("%s/", path.Join("/", rp))
}
