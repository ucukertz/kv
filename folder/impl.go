// Package folder implements key value store on local folder
package folder

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ucukertz/kv"
)

type Store[V []byte] struct {
	dir string
}

var _ kv.Store[[]byte] = (*Store[[]byte])(nil)
var _ kv.Bstore = (*Store[[]byte])(nil)

func Make(dir string) (*Store[[]byte], error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if errors.Is(err, os.ErrPermission) {
		return &Store[[]byte]{}, fmt.Errorf("%w folder make: %w", kv.ErrUnauthorized, err)
	} else if err != nil {
		return &Store[[]byte]{}, fmt.Errorf("%w folder make: %w", kv.ErrHalt, err)
	}
	return &Store[[]byte]{dir: dir}, nil
}

func (s *Store[V]) Set(k string, v []byte) error {
	dir := path.Join(s.dir, k)
	file, err := os.OpenFile(dir, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("%w folder set %s: %w", kv.ErrUnauthorized, k, err)
	} else if err != nil {
		return fmt.Errorf("%w folder set %s: %w", kv.ErrHalt, k, err)
	}
	defer file.Close()
	_, err = file.Write(v)
	if err != nil {
		return fmt.Errorf("%w folder set %s: %w", kv.ErrHalt, k, err)
	}
	return nil
}

func (s *Store[V]) Get(k string) ([]byte, error) {
	dir := path.Join(s.dir, k)
	file, err := os.Open(dir)
	if errors.Is(err, os.ErrNotExist) {
		return []byte{}, fmt.Errorf("%w folder get %s: %w", kv.ErrNotFound, k, err)
	} else if errors.Is(err, os.ErrPermission) {
		return []byte{}, fmt.Errorf("%w folder get %s: %w", kv.ErrUnauthorized, k, err)
	} else if err != nil {
		return []byte{}, fmt.Errorf("%w folder get %s: %w", kv.ErrHalt, k, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, fmt.Errorf("%w folder get %s: %w", kv.ErrHalt, k, err)
	}
	return content, nil
}

func (s *Store[V]) Delete(k string) error {
	dir := path.Join(s.dir, k)
	err := os.Remove(dir)
	if err != nil {
		return fmt.Errorf("%w folder del %s: %w", kv.ErrHalt, k, err)
	}
	return nil
}

func (s *Store[V]) Clear() error {
	err := os.RemoveAll(s.dir)
	if err != nil {
		return fmt.Errorf("%w folder clr: %w", kv.ErrHalt, err)
	}
	err = os.Mkdir(s.dir, os.ModePerm)
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("%w folder clr: %w", kv.ErrUnauthorized, err)
	} else if err != nil {
		return fmt.Errorf("%w folder clr: %w", kv.ErrHalt, err)
	}
	return nil
}

func (s *Store[V]) Purge() error {
	err := os.RemoveAll(s.dir)
	if err != nil {
		return fmt.Errorf("%w folder prg: %w", kv.ErrHalt, err)
	}
	return nil
}
