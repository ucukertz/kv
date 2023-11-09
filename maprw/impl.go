// Package maprw implements key value store with go map and RWMutex for concurrency
package maprw

import (
	"fmt"
	"sync"

	"github.com/ucukertz/kv"
)

type Store[V any] struct {
	m   map[string]V
	mtx sync.RWMutex
}

var _ kv.Store[any] = (*Store[any])(nil)
var _ kv.Bstore = (*Store[[]byte])(nil)

func Create[V any]() *Store[V] {
	return &Store[V]{m: map[string]V{}}
}

func (s *Store[V]) Set(k string, v V) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[k] = v
	return nil
}

func (s *Store[V]) Get(k string) (V, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	v, ok := s.m[k]
	var err error = nil
	if !ok {
		err = fmt.Errorf("maprw -> Reading key %s failed", k)
	}
	return v, err
}

func (s *Store[V]) Delete(k string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.m, k)
	return nil
}

func (s *Store[V]) Clear() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	clear(s.m)
	return nil
}

func (s *Store[V]) Close() error {
	s = Create[V]()
	return nil
}
