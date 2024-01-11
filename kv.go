// Package kv implements interface for key value store
package kv

// Minimal key value store interface
type Store[V any] interface {
	Set(k string, v V) error
	Get(k string) (V, error)
	Delete(k string) error
	Clear() error
	Purge() error
}

// Key value store with []byte as value
type Bstore Store[[]byte]

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

var (
	ErrHalt         = &Error{"KV halted"}
	ErrUnauthorized = &Error{"KV unauthorized"}
	ErrUnreachable  = &Error{"KV unreachable"}
	ErrNotFound     = &Error{"KV not found"}
)
