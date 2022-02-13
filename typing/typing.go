package typing

import (
	"errors"
	"reflect"
	"time"
)

func IntValue(v interface{}) (int, error) {
	switch v := v.(type) {
	case int:
		return v, nil

	case int32:
		return int(v), nil

	case uint32:
		return int(v), nil

	default:
		return 0, errors.New("parse error")
	}
}

func Int64Value(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int:
		return int64(v), nil

	case int32:
		return int64(v), nil

	case uint32:
		return int64(v), nil

	case int64:
		return v, nil

	case uint64:
		return int64(v), nil

	default:
		return 0, errors.New("parse error")
	}
}

func Float32Value(v interface{}) (float32, error) {
	switch v := v.(type) {
	case float32:
		return v, nil

	default:
		return 0.0, errors.New("parse error")
	}
}

func Float64Value(v interface{}) (float64, error) {
	switch v := v.(type) {
	case float32:
		return float64(v), nil

	case float64:
		return v, nil

	default:
		return 0.0, errors.New("parse error")
	}
}

func BoolValue(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil

	default:
		return false, errors.New("parse error")
	}
}

func StringValue(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil

	default:
		return "", errors.New("parse error")
	}
}

func BytesValue(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil

	default:
		return nil, errors.New("parse error")
	}
}

func TimeValue(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case time.Time:
		return v, nil

	default:
		return time.Now(), errors.New("parse error")
	}
}

func NilValue(v interface{}) (interface{}, error) {
	if v == nil || reflect.ValueOf(v).IsNil() {
		return nil, nil
	} else {
		return nil, errors.New("parse error")
	}
}
