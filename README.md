# Badger
Simple Router multiplexer for web api's based on [httprouter](https://github.com/julienschmidt/httprouter) adding the feature to add Middlewares to specific group of routes.

## Why a new Router?
After looking into another options like [gorilla/mux](https://github.com/gorilla/mux) and [Gin](https://github.com/gin-gonic/gin) I realized both had advantages and restrictions.

Gorila has plenty of middlewares that respect the [http.Handler](https://golang.org/pkg/net/http/#Handler) interface, which makes easy to find and develop new middlewares and custom handlers, but it lacks context as parameter.

### But they have Context before it was cool
Yes, they have a [context](https://github.com/gorilla/context), but this context relies on a global variable shared between all routines which causes unnecessary concurrency.

#### Route declaration
Really, in what world this:

``` golang
r.HandleFunc("/products", ProductsHandler).
  Methods("GET")
```

is easier to use than this:

``` golang
r.GET("/products", func(c *gin.Context) {
	// DO work
})
```

Kidding aside, Gorilla has plenty more of built in functions regarding parameter validations within mux, url scheme matching and other things, which will be not considered here as it is not a requirement.

### Ok, then Gin's context parameters is the solution and route grouping is the solution
Yes, but with some drawbacks.

Aside being based on [httprouter](https://github.com/julienschmidt/httprouter) which makes [Gin](https://github.com/gin-gonic/gin) 30x faster than [gorilla/mux](https://github.com/gorilla/mux) it has a route grouping which you can configure some routes to have some middlewares.

#### Perfect! Wait... not so fast
Gin's implementation of context is pretty close to what I was looking for, it [returns a pointer Context object](https://github.com/gin-gonic/gin/blob/master/gin.go#L320) using [sync.Pool.Get()](https://golang.org/pkg/sync/#Pool) method and then the interface to handle request uses only this context, which is a all in one object.

It contains custom ResponseWriter interface, custom query params interface, JSON writer, and many other funcionalities, it is almost magical.

## Now what?
How can we get the best of the 2 worlds?

The requirements are simple:
* Be compliant to [Golang/net](https://golang.org/pkg/net/http) interfaces to be easier to find new middlewares (and use gorilla's)
* Have Context parameter (not being global or magical)

# The challenge
Build a mux that implements the requirements above

## How does it work?
#### Yes it is a whole server... I'll work on more simple examples later (promisse)

``` golang
package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/hugoluchessi/badger"
	"go.uber.org/zap"
)

func main() {
	// Create new Mux
	mux := badger.NewMux()

	// Create new router group
	router1 := mux.AddRouter("v1")

	// Adds an handler for route someget
	router1.Get("someget", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello, I'm router1, accessed by url %s.", req.URL.Path[1:])
	}))

	// Create another router group
	router2 := mux.AddRouter("v2")

	// Adds an handler for route someget
	router2.Get("someget", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello, I'm router2, also accessed by url %s, and using transaction id %d.", req.URL.Path[1:], req.Context().Value("TransactionId"))
	}))

	// Example logger uber-zap(https://github.com/uber-go/zap)
	logger := zap.NewExample()
	defer logger.Sync()

	// Define a middleware used by router 1
	router1.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Info("Started")
			h.ServeHTTP(rw, req)
			logger.Info("Finished")
		})
	})

	// Define another middleware used by router 1
	router2.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			// This number could be a header value
			ctx = context.WithValue(ctx, "TransactionId", rand.Uint64())
			req = req.WithContext(ctx)
			h.ServeHTTP(rw, req)
		})
	})

	// Define a middleware used by router 2
	router2.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Info("Started middleware 2")
			h.ServeHTTP(rw, req)
			logger.Info("Finished middleware 2")
		})
	})

	http.ListenAndServe(":8080", mux)
}

```

On the terminal, use curl to access the routes:

```
curl localhost:8080/v1/someget   
Hello, I'm router1, accessed by url v1/someget.

curl localhost:8080/v2/someget
=> Hello, I'm router2, also accessed by url v2/someget, and using transaction id 5577006791947779410.
```

Yes! It works, and the transaction ID correctly arrived the handler function. Also there is log messages:

```
{"level":"info","msg":"Started"}
{"level":"info","msg":"Finished"}
{"level":"info","msg":"Started middleware 2"}
{"level":"info","msg":"Finished middleware 2"}
```
