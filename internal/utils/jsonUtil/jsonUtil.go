package jsonUtil

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
)

func MarshalValue[T any](val T) (driver.Value, error) {
	if isZero(val) {
		return nil, nil
	}
	return json.Marshal(val)
}

func UnmarshalValue[T any](value interface{}, dest *T) error {
	if value == nil {
		*dest = *new(T) // 设置为类型 T 的零值
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("jsonUtil: expected []byte from database")
	}
	return json.Unmarshal(bytes, dest)
}

func isZero[T any](v T) bool {
	return reflect.ValueOf(v).IsZero()
}
