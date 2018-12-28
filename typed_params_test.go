package badger_test

import (
	"testing"

	"github.com/hugoluchessi/badger"
)

func TestGetString(t *testing.T) {
	key := "map"
	value := "mapvalue!"

	dict := map[string]string{key: value}

	typedmap := badger.CreateTypedParams(dict)

	rvalue, err := typedmap.GetString(key)

	if err != nil {
		t.Error("Test failed, err must be nil.")
	}

	if rvalue == "" {
		t.Errorf("Test failed, expected value to be '%s' got '%s'.", value, rvalue)
	}
}

func TestGetInexistingStringKey(t *testing.T) {
	key := "map"
	value := "keyvalue!"

	dict := map[string]string{key: value}

	typedmap := badger.CreateTypedParams(dict)

	rvalue, err := typedmap.GetString("someotherkey")

	if err == nil {
		t.Error("Test failed, err must not be nil.")
	}

	if rvalue != "" {
		t.Errorf("Test failed, expected value to be '%s' got '%s'.", "", rvalue)
	}
}

func TestGetint(t *testing.T) {
	key := "map"
	value := "1223"

	dict := map[string]string{key: value}

	typedmap := badger.CreateTypedParams(dict)

	rvalue, err := typedmap.GetInt(key)

	if err != nil {
		t.Error("Test failed, err must be nil.")
	}

	if rvalue == 0 {
		t.Errorf("Test failed, expected value to be '%s' got '%d'.", value, rvalue)
	}
}

func TestGetInexistingIntKey(t *testing.T) {
	key := "map"
	value := "1223"

	dict := map[string]string{key: value}

	typedmap := badger.CreateTypedParams(dict)

	rvalue, err := typedmap.GetInt("someotherkey")

	if err == nil {
		t.Error("Test failed, err must not be nil.")
	}

	if rvalue != 0 {
		t.Errorf("Test failed, expected value to be '%d' got '%d'.", 0, rvalue)
	}
}
