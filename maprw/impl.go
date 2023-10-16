// Package maprw implements key value store with go map and RWMutex for concurrency
package maprw

import (
	"fmt"
	"sync"
)

type Store[V any] struct {
	m    map[string]V
	lock *sync.RWMutex
}

func Create[V any]() *Store[V] {
	return &Store[V]{map[string]V{}, &sync.RWMutex{}}
}

func (s *Store[V]) Set(k string, v V) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[k] = v
	return nil
}

func (s *Store[V]) Get(k string) (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	v, ok := s.m[k]
	var err error = nil
	if !ok {
		err = fmt.Errorf("Reading key %s failed", k)
	}
	return v, err
}

func (s *Store[V]) Delete(k string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, k)
	return nil
}

func (s *Store[V]) Clear() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for k := range s.m {
		delete(s.m, k)
	}
	return nil
}

func (s *Store[V]) Close() error {
	s = Create[V]()
	return nil
}
