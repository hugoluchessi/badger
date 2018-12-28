package badger_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hugoluchessi/badger"
	"github.com/julienschmidt/httprouter"
)

func TestCreateRouteParams(t *testing.T) {
	key := "map"
	value := "mapvalue!"

	params := httprouter.Params{}
	params = append(params, httprouter.Param{Key: key, Value: value})

	typedmap := badger.CreateRouteParams(params)

	rvalue, err := typedmap.GetString(key)

	if err != nil {
		t.Error("Test failed, err must be nil.")
	}

	if rvalue == "" {
		t.Errorf("Test failed, expected value to be '%s' got '%s'.", value, rvalue)
	}
}

func TestCreateRouteParamsWithNoParams(t *testing.T) {
	key := "map"

	params := httprouter.Params{}
	params = append(params, httprouter.Param{})

	typedmap := badger.CreateRouteParams(params)

	rvalue, err := typedmap.GetString(key)

	if err == nil {
		t.Error("Test failed, err must not be nil.")
	}

	if rvalue != "" {
		t.Errorf("Test failed, expected value to be '%s' got '%s'.", "", rvalue)
	}
}

func TestCreateRouteIntParams(t *testing.T) {
	key := "map"
	value := "1234"

	params := httprouter.Params{}
	params = append(params, httprouter.Param{Key: key, Value: value})

	typedmap := badger.CreateRouteParams(params)

	rvalue, err := typedmap.GetInt(key)

	if err != nil {
		t.Error("Test failed, err must be nil.")
	}

	if rvalue == 0 {
		t.Errorf("Test failed, expected value to be '%s' got '%d'.", value, rvalue)
	}
}

func TestGetRouteParamsFromRequest(t *testing.T) {
	key := "name"
	value := "cool"

	dict := map[string]string{
		key: value,
	}

	typed := badger.CreateTypedParams(dict)
	req, _ := http.NewRequest("GET", "nowhere", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, badger.RouteParamsKey, typed)
	req = req.WithContext(ctx)

	dictexpected := badger.GetRouteParamsFromRequest(req)

	if rvalue, _ := dictexpected.GetString(key); rvalue == "" {
		t.Errorf("Test failed, expected value to be '%s' got '%s'.", value, rvalue)
	}
}
