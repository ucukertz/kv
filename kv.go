// Package kv implements interface for key value store
package kv

// Minimal key value store interface
type Store[V any] interface {
	Set(k string, v V) error
	Get(k string) (V, error)
	Delete(k string) error
	Clear() error
	Close() error
}

// Key value store with []byte as value
type Bstore Store[[]byte]
