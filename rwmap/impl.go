// Package rwmap implements key value store with go map and RWMutex for concurrency
package rwmap

import (
	"fmt"
	"sync"

	"github.com/ucukertz/kv"
)

type Store[V any] struct {
	m map[string]V
	*sync.RWMutex
}

var _ kv.Store[any] = (*Store[any])(nil)
var _ kv.Bstore = (*Store[[]byte])(nil)

func Make[V any]() *Store[V] {
	return &Store[V]{m: map[string]V{}}
}

func (s *Store[V]) Set(k string, v V) error {
	s.Lock()
	defer s.Unlock()
	s.m[k] = v
	return nil
}

func (s *Store[V]) Get(k string) (V, error) {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.m[k]
	var err error
	if !ok {
		err = fmt.Errorf("%w rwmap get %s: %w", kv.ErrNotFound, k, err)
	}
	return v, err
}

func (s *Store[V]) Delete(k string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.m, k)
	return nil
}

func (s *Store[V]) Clear() error {
	s.Lock()
	defer s.Unlock()
	clear(s.m)
	return nil
}

func (s *Store[V]) Purge() error {
	s = nil
	return nil
}
