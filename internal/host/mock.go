package host

import (
	"context"
	"fmt"
)

type Mock struct {
	Results map[string]Value
	Calls   []string
}

func (m *Mock) Call(_ context.Context, name string, _ []Value) (Value, error) {
	m.Calls = append(m.Calls, name)
	v, ok := m.Results[name]
	if !ok {
		return Value{}, fmt.Errorf("mock host has no result for %s", name)
	}
	return v, nil
}
