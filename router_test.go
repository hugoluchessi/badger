package badger

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)

	if router == nil {
		t.Error("Test failed, router should not be nil")
	}
}

func AssertHandlerFunc(ehandlerheaderkey string, ehandlerheadervalue string) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add(ehandlerheaderkey, ehandlerheadervalue)
	}
}

func AssertRoute(t *testing.T, route Route, emethod string, epath string, ehandlerheaderkey string, ehandlerheadervalue string) {
	req, _ := http.NewRequest(emethod, epath, nil)
	res := httptest.NewRecorder()

	if route.method != emethod {
		t.Errorf("Test failed, wrong route method, got '%s' expected '%s'.", route.method, emethod)
	}

	if route.path != epath {
		t.Errorf("Test failed, wrong route path, got '%s' expected '%s'.", route.path, epath)
	}

	route.handler.ServeHTTP(res, req)

	handlervalue := res.Header().Get(ehandlerheaderkey)

	if handlervalue != ehandlerheadervalue {
		t.Errorf("Test failed, wrong route handler value, got '%s' expected '%s'.", handlervalue, ehandlerheadervalue)
	}
}

func TestDelete(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "DELETE"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Delete(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestGet(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "GET"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Get(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestHead(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "HEAD"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Head(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestOptions(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "OPTIONS"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Options(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestPatch(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "PATCH"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Patch(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestPost(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "POST"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Post(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestPut(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "PUT"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Put(path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestHandle(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "SOMEMETHOD"
	path := "/somePath"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Handle(method, path, handler)

	if len(router.routes) == 0 {
		t.Error("Test failed, there should be one route on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
}

func TestGetTwice(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	method := "GET"
	path := "/somePath"
	path2 := "/somePath2"
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Get(path, handler)
	router.Get(path2, handler)

	if len(router.routes) < 2 {
		t.Error("Test failed, there should be two routes on router.")
	}

	AssertRoute(t, router.routes[0], method, path, handlerheaderkey, handlerheadervalue)
	AssertRoute(t, router.routes[1], method, path2, handlerheaderkey, handlerheadervalue)
}

func TestMultipleHandlers(t *testing.T) {
	basepath := ""
	router := NewRouter(basepath)
	handlerheaderkey := "some key"
	handlerheadervalue := "some value"

	method1 := "GET"
	path1 := "/somePath"

	method2 := "POST"
	path2 := "/somePath2"

	method3 := "OTHER"
	path3 := "/somePath3"

	handler := http.HandlerFunc(AssertHandlerFunc(handlerheaderkey, handlerheadervalue))

	router.Get(path1, handler)
	router.Post(path2, handler)
	router.Handle(method3, path3, handler)

	if len(router.routes) < 3 {
		t.Error("Test failed, there should be three routes on router.")
	}

	AssertRoute(t, router.routes[0], method1, path1, handlerheaderkey, handlerheadervalue)
	AssertRoute(t, router.routes[1], method2, path2, handlerheaderkey, handlerheadervalue)
	AssertRoute(t, router.routes[2], method3, path3, handlerheaderkey, handlerheadervalue)
}

func MyTestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("X-Middleware", "YEAH!")
		h.ServeHTTP(res, req)
		res.Header().Add("X-Middleware-2", "YEAH!")
	})
}

func MyTestMiddleware2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("X-Middleware2", "YEAH!")
		h.ServeHTTP(res, req)
		res.Header().Add("X-Middleware2-2", "YEAH!")
	})
}

func TestUse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health-check", nil)
	res := httptest.NewRecorder()

	basepath := ""
	router := NewRouter(basepath)

	router.Use(MyTestMiddleware)

	if len(router.middlewares) < 1 {
		t.Error("Test failed, there should be one middleware on router.")
	}

	h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		header1 := res.Header().Get("X-Middleware")
		header2 := res.Header().Get("X-Middleware-2")

		if header1 != "YEAH!" {
			t.Errorf("Test failed, invalid 'X-Middleware' header value, got '%s' expected '%s'.", header1, "YEAH!")
		}

		if header2 != "" {
			t.Errorf("Test failed, invalid 'X-Middleware-2' header value, got '%s' expected '%s'.", header2, "")
		}
	})

	builthandler := router.middlewares[0].BuildHandler(h)

	builthandler.ServeHTTP(res, req)

	header2 := res.Header().Get("X-Middleware-2")

	if header2 != "YEAH!" {
		t.Errorf("Test failed, invalid 'X-Middleware' header value, got '%s' expected '%s'.", header2, "YEAH!")
	}
}

func TestUseMultiple(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health-check", nil)
	res := httptest.NewRecorder()

	basepath := ""
	router := NewRouter(basepath)

	router.Use(MyTestMiddleware)
	router.Use(MyTestMiddleware2)

	if len(router.middlewares) < 2 {
		t.Error("Test failed, there should be two middlewares on router.")
	}

	h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		mw1header1 := res.Header().Get("X-Middleware")
		mw2header1 := res.Header().Get("X-Middleware2")
		mw1header2 := res.Header().Get("X-Middleware-2")
		mw2header2 := res.Header().Get("X-Middleware2-2")

		if mw1header1 != "YEAH!" {
			t.Errorf("Test failed, invalid 'X-Middleware' header value, got '%s' expected '%s'.", mw1header1, "YEAH!")
		}

		if mw1header2 != "" {
			t.Errorf("Test failed, invalid 'X-Middleware-2' header value, got '%s' expected '%s'.", mw1header2, "")
		}

		if mw2header1 != "YEAH!" {
			t.Errorf("Test failed, invalid 'X-Middleware' header value, got '%s' expected '%s'.", mw2header1, "YEAH!")
		}

		if mw2header2 != "" {
			t.Errorf("Test failed, invalid 'X-Middleware-2' header value, got '%s' expected '%s'.", mw2header2, "")
		}
	})

	builthandler := router.middlewares[0].BuildHandler(h)
	builthandler = router.middlewares[1].BuildHandler(builthandler)

	builthandler.ServeHTTP(res, req)

	mw1header2 := res.Header().Get("X-Middleware-2")
	mw2header2 := res.Header().Get("X-Middleware2-2")

	if mw1header2 != "YEAH!" {
		t.Errorf("Test failed, invalid 'X-Middleware' header value, got '%s' expected '%s'.", mw1header2, "YEAH!")
	}

	if mw2header2 != "YEAH!" {
		t.Errorf("Test failed, invalid 'X-Middleware' header value, got '%s' expected '%s'.", mw2header2, "YEAH!")
	}
}
