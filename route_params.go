package badger

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type routeParamsKey struct{}

var RouteParamsKey = routeParamsKey{}

type RouteParams struct {
	params TypedParams
}

func CreateRouteParams(rps httprouter.Params) TypedParams {
	dict := make(map[string]string)

	for _, rp := range rps {
		if rp.Key == "" {
			continue
		}

		dict[rp.Key] = rp.Value
	}

	return CreateTypedParams(dict)
}

func GetRouteParamsFromRequest(req *http.Request) TypedParams {
	ctx := req.Context()
	return ctx.Value(RouteParamsKey).(TypedParams)
}
