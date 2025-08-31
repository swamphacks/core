package web

import (
	"errors"
	"net/url"
	"strconv"
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
