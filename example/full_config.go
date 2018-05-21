package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/hugoluchessi/badger"
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
	router2.Get("someget/:someparam", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		routeparams := badger.GetRouteParamsFromRequest(req)
		someparam, err := routeparams.GetString("someparam")

		if err != nil {
			fmt.Fprintf(res, "Shoot, param 'someparam' had some problems. Error: %s", err.Error())
		}

		fmt.Fprintf(
			res,
			"Hello, I'm router2, also accessed by url %s, and using transaction id %d, and route param %s",
			req.URL.Path[1:],
			req.Context().Value("TransactionId"),
			someparam,
		)
	}))

	// Example logger uber-zap(https://github.com/uber-go/zap)
	logger := log.New(os.Stdout, "", log.Flags())

	// Define a middleware used by router 1
	router1.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Print("Started")
			h.ServeHTTP(rw, req)
			logger.Print("Finished")
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
			logger.Print("Comecou middleware 2")
			h.ServeHTTP(rw, req)
			logger.Print("Terminou middleware 2")
		})
	})

	// Method not allowed for existing route handler
	mux.MethodNotAllowed = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger.Print("Can't do with this method")
	})

	// Route not found handler
	mux.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		logger.Print("Sorry, wrong route")
	})

	// Panic handler
	mux.PanicHandler = func(rw http.ResponseWriter, req *http.Request, p interface{}) {
		logger.Panic(fmt.Sprintf("Panicked with '%s'", p.(string)))
	}

	http.ListenAndServe(":8080", mux)
}
