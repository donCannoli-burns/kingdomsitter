package host

import "context"

type Value struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// Host is the only runtime boundary allowed to answer real KoLmafia-backed
// calls. Parser/analyzer/agent packages do not receive direct game authority.
type Host interface {
	Call(ctx context.Context, name string, args []Value) (Value, error)
}
