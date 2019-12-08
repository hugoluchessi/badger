[![CircleCI](https://circleci.com/gh/hugoluchessi/badger/tree/master.svg?style=shield)](https://circleci.com/gh/hugoluchessi/badger/tree/master)
# Badger
Simple Route multiplexer for web api's based on [httprouter](https://github.com/julienschmidt/httprouter) adding the feature to add Middlewares to specific group of routes.

## Why a new Router?
Please see [the reason here](https://gist.github.com/hugoluchessi/db89f6f0fae0aced6251153bb97ee485).

## Features
* Fast and versatile routing
* Easy to use route grouping
* Middlewares
* 100% stdlib interfaces

## How does it work?
### Simple server

``` golang
package main

import (
	"net/http"

	"github.com/hugoluchessi/badger"
)

func main() {
	// Create new Mux
	mux := badger.NewMux()

	// Create new router group
	router := mux.AddRouter("v1")

	// Handler got GET Products
	router.Get("products", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		products := YourProductsDAL.all
		fmt.Fprintf(res, toJson(products))
	}))

	// Handler got GET Products
	router.Get("products/:id", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rp := badger.GetRouteParamsFromRequest(req)
		id, _ := rp.GetInt("id")

		products := YourProductsDAL.where("id = %d", id)
		fmt.Fprintf(res, toJson(products))
	}))

	// Define a middleware used by router 1
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			logger.Info("Started")
			h.ServeHTTP(rw, req)
			logger.Info("Finished")
		})
	})

	http.ListenAndServe(":8080", mux)
}

```

## Performance
[Here](https://github.com/hugoluchessi/go-http-routing-benchmark) is the project with the benchmark.

| Test | No of Operations | Time by op | Bytes per Op   | Allocations per op |
|:-------------|----------:|-----------:|----------:|----------:|
BenchmarkBadger_Param|1000000|2514 ns/op|944 B/op|13 allocs/op|
BenchmarkGin_Param|20000000|61.7 ns/op|0 B/op|0 allocs/op|
BenchmarkGorillaMux_Param|500000|3263 ns/op|1280 B/op|10 allocs/op|
BenchmarkHttpRouter_Param|20000000|110 ns/op|32 B/op|1 allocs/op|
BenchmarkMartini_Param|300000|5567 ns/op|1072 B/op|10 allocs/op|
BenchmarkBadger_Param5|500000|3267 ns/op|1024 B/op|13 allocs/op|
BenchmarkGin_Param5|20000000|108 ns/op|0 B/op|0 allocs/op|
BenchmarkGorillaMux_Param5|300000|4732 ns/op|1344 B/op|10 allocs/op|
BenchmarkHttpRouter_Param5|5000000|321 ns/op|160 B/op|1 allocs/op|
BenchmarkMartini_Param5|200000|6797 ns/op|1232 B/op|11 allocs/op|
BenchmarkAce_Param20|1000000|1581 ns/op|640 B/op|1 allocs/op|
BenchmarkBadger_Param20|500000|3768 ns/op|1104 B/op|13 allocs/op|
BenchmarkGin_Param20|5000000|277 ns/op|0 B/op|0 allocs/op|
BenchmarkGorillaMux_Param20|100000|12494 ns/op|3452 B/op|12 allocs/op|
BenchmarkHttpRouter_Param20|1000000|1528 ns/op|640 B/op|1 allocs/op|
BenchmarkMartini_Param20|100000|12697 ns/op|3596 B/op|13 allocs/op|
BenchmarkBadger_ParamWrite|1000000|2866 ns/op|944 B/op|13 allocs/op|
BenchmarkGin_ParamWrite|10000000|163 ns/op|0 B/op|0 allocs/op|
BenchmarkGorillaMux_ParamWrite|500000|3246 ns/op|1280 B/op|10 allocs/op|
BenchmarkHttpRouter_ParamWrite|10000000|135 ns/op|32 B/op|1 allocs/op|
BenchmarkMartini_ParamWrite|200000|6270 ns/op|1176 B/op|14 allocs/op|
