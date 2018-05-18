package badger

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

const GET = "GET"
const POST = "POST"
const RouterBasePath1 = "v1"
const RouterBasePath2 = "v2"
const RoutePath1 = "somepath"
const RoutePath2 = "otherpath"
const RouteHeaderKey1 = "X-Route"
const RouteHeaderKey2 = "X-Route-2"
const MiddlewareHeaderKey1 = "X-Middleware"
const MiddlewareHeaderKey2 = "X-Middleware-2"
const HeadersExpectedValue = "YEAH!"

func AssertHeader(t *testing.T, res http.ResponseWriter, key string, expectedvalue string) {
	value := res.Header().Get(key)

	if value != expectedvalue {
		t.Errorf("Test failed, invalid '%s' header value, got '%s' expected '%s'.", key, value, expectedvalue)
	}
}

func TestNewMux(t *testing.T) {
	mux := NewMux()

	if mux == nil {
		t.Error("Test failed, mux must not be nil")
	}
}

func TestAddRouter(t *testing.T) {
	path := "root"
	mux := NewMux()
	router := mux.AddRouter(path)

	if len(mux.routers) != 1 {
		t.Error("Test failed, there must be one router on mux.")
	}

	if router == nil {
		t.Error("Test failed, router must not be nil")
	}
}

func TestAddMutlipleRouters(t *testing.T) {
	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)
	router2 := mux.AddRouter(RouterBasePath2)

	if len(mux.routers) != 2 {
		t.Error("Test failed, there must be two routers on mux.")
	}

	if router == nil {
		t.Error("Test failed, router must not be nil")
	}

	if router2 == nil {
		t.Error("Test failed, router2 must not be nil")
	}
}

func TestServeHTTPNoRouters(t *testing.T) {
	req, _ := http.NewRequest(GET, RouterBasePath1, nil)
	res := httptest.NewRecorder()

	mux := NewMux()

	mux.ServeHTTP(res, req)
}

func TestServeHTTPOneRouterNoRoutesNoMiddlewares(t *testing.T) {
	req, _ := http.NewRequest(GET, RouterBasePath1, nil)
	res := httptest.NewRecorder()

	mux := NewMux()
	_ = mux.AddRouter(RouterBasePath1)

	mux.ServeHTTP(res, req)
}

func TestServeHTTPOneRouterOneRouteNoMiddlewares(t *testing.T) {
	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()

	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)
	router.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	mux.ServeHTTP(res, req)

	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)
}

func TestServeHTTPOneRouterOneRouteOneMiddleware(t *testing.T) {
	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()

	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)

	router.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey1, HeadersExpectedValue)
		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(MiddlewareHeaderKey1, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	mux.ServeHTTP(res, req)

	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)
}

func TestServeHTTPOneRouterOneRouteTwoMiddlewares(t *testing.T) {
	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()

	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)

	router.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey1, HeadersExpectedValue)
		AssertHeader(t, rw, MiddlewareHeaderKey2, HeadersExpectedValue)

		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			AssertHeader(t, rw, MiddlewareHeaderKey2, HeadersExpectedValue)

			rw.Header().Add(MiddlewareHeaderKey1, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(MiddlewareHeaderKey2, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	mux.ServeHTTP(res, req)

	AssertHeader(t, res, MiddlewareHeaderKey1, HeadersExpectedValue)
	AssertHeader(t, res, MiddlewareHeaderKey2, HeadersExpectedValue)
	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)
}

func TestServeHTTPOneRouterTwoRouteNoMiddleware(t *testing.T) {
	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)

	router.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	router.Post(RoutePath2, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add(RouteHeaderKey2, HeadersExpectedValue)
	}))

	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)

	req, _ = http.NewRequest(POST, path.Join("/", RouterBasePath1, RoutePath2), nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, RouteHeaderKey2, HeadersExpectedValue)
}

func FuncWithMid2(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add(RouteHeaderKey2, HeadersExpectedValue)
}

func TestServeHTTPOneRouterTwoRouteOneMiddleware(t *testing.T) {
	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)

	router.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey1, HeadersExpectedValue)
		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	router.Get(RoutePath2, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey1, HeadersExpectedValue)
		rw.Header().Add(RouteHeaderKey2, HeadersExpectedValue)
	}))

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(MiddlewareHeaderKey1, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)

	req, _ = http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath2), nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, MiddlewareHeaderKey1, HeadersExpectedValue)
	AssertHeader(t, res, RouteHeaderKey2, HeadersExpectedValue)
}

func TestServeHTTPOneRouterOneRouteTwoMiddleware(t *testing.T) {
	mux := NewMux()
	router := mux.AddRouter(RouterBasePath1)

	router.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey1, HeadersExpectedValue)
		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			AssertHeader(t, rw, MiddlewareHeaderKey2, HeadersExpectedValue)
			rw.Header().Add(MiddlewareHeaderKey1, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(MiddlewareHeaderKey2, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, MiddlewareHeaderKey1, HeadersExpectedValue)
	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)
}

func TestServeHTTPTwoRoutersOneRouteOneMiddleware(t *testing.T) {
	mux := NewMux()
	router1 := mux.AddRouter(RouterBasePath1)
	router2 := mux.AddRouter(RouterBasePath2)

	router1.Get(RoutePath1, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey1, HeadersExpectedValue)
		rw.Header().Add(RouteHeaderKey1, HeadersExpectedValue)
	}))

	router1.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(MiddlewareHeaderKey1, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	router2.Get(RoutePath2, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		AssertHeader(t, rw, MiddlewareHeaderKey2, HeadersExpectedValue)
		rw.Header().Add(RouteHeaderKey2, HeadersExpectedValue)
	}))

	router2.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(MiddlewareHeaderKey2, HeadersExpectedValue)
			h.ServeHTTP(rw, req)
		})
	})

	req, _ := http.NewRequest(GET, path.Join("/", RouterBasePath1, RoutePath1), nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, MiddlewareHeaderKey1, HeadersExpectedValue)
	AssertHeader(t, res, RouteHeaderKey1, HeadersExpectedValue)
	AssertHeader(t, res, MiddlewareHeaderKey2, "")
	AssertHeader(t, res, RouteHeaderKey2, "")

	req, _ = http.NewRequest(GET, path.Join("/", RouterBasePath2, RoutePath2), nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	AssertHeader(t, res, MiddlewareHeaderKey1, "")
	AssertHeader(t, res, RouteHeaderKey1, "")
	AssertHeader(t, res, MiddlewareHeaderKey2, HeadersExpectedValue)
	AssertHeader(t, res, RouteHeaderKey2, HeadersExpectedValue)
}

func TestServeHTTPWithParams(t *testing.T) {
	mux := NewMux()
	routepath := path.Join("/", RouterBasePath1, RoutePath1)
	router := mux.AddRouter("")

	router.Get(path.Join(routepath, "/:testparam"), http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		params := GetRouteParamsFromRequest(r)
		testparam, _ := params.GetString("testparam")

		if testparam != "test" {
			t.Errorf("Test failed, invalid param, expected '%s' got '%s'.", testparam, "test")
		}
	}))

	req, _ := http.NewRequest(GET, path.Join(routepath, "/test", "?test=param1&test2=param2"), nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)
}
