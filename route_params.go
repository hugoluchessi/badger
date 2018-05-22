package badger

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type routeParamsKey struct{}

// RouteParamsKey is the key to find route params in context
var RouteParamsKey = routeParamsKey{}

// RouteParams are the params found in route named parameters
// they are a TypedParams type, which allows to retrieve typed info
type RouteParams struct {
	params TypedParams
}

// CreateRouteParams converts httprouter.Params to TypedParams
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

// GetRouteParamsFromRequest retrieves Typed Route params from given request
func GetRouteParamsFromRequest(req *http.Request) TypedParams {
	ctx := req.Context()
	return ctx.Value(RouteParamsKey).(TypedParams)
}
