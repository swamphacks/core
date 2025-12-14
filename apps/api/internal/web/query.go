package web

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrMalformedField = errors.New("malformed field")
)

// ParseParamBoolean parses a boolean query parameter from the given url.Values.
// If the parameter is missing or cannot be parsed as a boolean, it returns defaultVal.
//
// Parameters:
//   - queryParams: the URL query parameters to read from
//   - key: the name of the query parameter
//   - defaultVal: pointer to the default boolean value to use if parsing fails or the parameter is absent
//
// Returns:
//   - pointer to the parsed boolean value, or defaultVal if missing or invalid
//   - error if field was malformed/not boolean
func ParseParamBoolean(queryParams url.Values, key string, defaultVal *bool) (*bool, error) {
	v := queryParams.Get(key)
	if v == "" {
		return defaultVal, nil
	}

	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return nil, ErrMalformedField
	}

	return &parsed, nil
}

func ParseParamInt32(queryParams url.Values, key string, defaultVal *int32) (*int32, error) {
	v := queryParams.Get(key)
	if v == "" {
		return defaultVal, nil
	}

	parsed, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return nil, ErrMalformedField
	}

	val := int32(parsed)

	return &val, nil
}

func ParseParamString(queryParams url.Values, key string, defaultVal *string) *string {
	v := queryParams.Get(key)
	if v == "" {
		return defaultVal
	}

	return &v
}

// This deprecates functions from above.
func ParseParam[T any](query url.Values, key string, def *T) (*T, error) {
	v := query.Get(key)
	if v == "" {
		return def, nil
	}

	var out T
	switch any(out).(type) {
	case bool:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return nil, ErrMalformedField
		}
		out = any(parsed).(T)

	case int32:
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, ErrMalformedField
		}
		out = any(int32(parsed)).(T)

	case string:
		out = any(v).(T)

	case uuid.UUID:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return nil, ErrMalformedField
		}
		out = any(parsed).(T)

	default:
		return nil, errors.New("unsupported type for ParseParam")
	}

	return &out, nil
}
