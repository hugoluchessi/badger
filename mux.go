package badger

import (
	"net/http"
	"path"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type HandleFunc func(http.Handler) httprouter.Handle

type Mux struct {
	routers    []*Router
	mainrouter *httprouter.Router
	lock       sync.RWMutex
}

func NewMux() *Mux {
	return &Mux{[]*Router{}, nil, sync.RWMutex{}}
}

func (mux *Mux) AddRouter(path string) *Router {
	mux.lock.Lock()
	defer mux.lock.Unlock()

	router := NewRouter(path)
	mux.routers = append(mux.routers, router)

	return router
}

func (mux *Mux) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

	for _, router := range mux.routers {
		routerroutes := router.buildRoutes()

		for _, route := range routerroutes {
			mux.mainrouter.Handle(
				route.method,
				route.path,
				// FIXME: Ignore params for now
				HandleFunc(func(h http.Handler) httprouter.Handle {
					return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
		// Ensure path starts with /
		p := path.Join("/", r.basepath, route.path)
		builtroute := Route{route.method, p, route.handler}

		for _, middleware := range r.middlewares {
			builtroute.handler = middleware.BuildHandler(builtroute.handler)
		}

		builtroutes = append(builtroutes, builtroute)
	}

	return builtroutes
}
