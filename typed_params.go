package badger

import (
	"fmt"
	"strconv"
)

const notFoundErrorMessageFormat = "Key '%s' could not be found."

// TypedParams is a helper struct fo handling param objects, has helper
// functions to retrieve typed data
type TypedParams struct {
	params map[string]string
}

// CreateTypedParams creates and returns TypedParams
func CreateTypedParams(params map[string]string) TypedParams {
	return TypedParams{params}
}

// GetString returns an string value for the given key, returns "" and an error
// in case key was not found
func (t TypedParams) GetString(key string) (string, error) {
	if val, ok := t.params[key]; ok {
		return val, nil
	}

	return "", fmt.Errorf(notFoundErrorMessageFormat, key)
}

// GetInt returns an integer value for the given key, returns 0 and an error
// in case key was not found and also returns an error in caso conversion
// is not successful
func (t TypedParams) GetInt(key string) (int, error) {
	value, err := t.GetString(key)

	if err != nil {
		return 0, err
	}

	ivalue, err := strconv.Atoi(value)

	return ivalue, err
}
