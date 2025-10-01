package parse

import "encoding/json"

type Optional[T any] struct {
	Value   T
	Present bool
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.Present = true
	return json.Unmarshal(data, &o.Value)
}
