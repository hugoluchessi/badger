package badger

import (
	"errors"
	"fmt"
	"strconv"
)

const notFoundErrorMessageFormat = "Key '%s' could not be found."

type TypedParams struct {
	params map[string]string
}

func CreateTypedParams(params map[string]string) TypedParams {
	return TypedParams{params}
}

func (t TypedParams) GetString(key string) (string, error) {
	if val, ok := t.params[key]; ok {
		return val, nil
	}

	return "", errors.New(fmt.Sprintf(notFoundErrorMessageFormat, key))
}

func (t TypedParams) GetInt(key string) (int, error) {
	value, err := t.GetString(key)

	if err != nil {
		return 0, err
	}

	ivalue, err := strconv.Atoi(value)

	return ivalue, err
}
