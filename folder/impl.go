// Package folder implements key value store on local folder
package folder

import (
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

func Create(dir string) (*Store[[]byte], error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return &Store[[]byte]{}, fmt.Errorf("folder -> %w", err)
	}
	return &Store[[]byte]{dir: dir}, nil
}

func (s *Store[V]) Set(k string, v []byte) error {
	dir := path.Join(s.dir, k)
	file, err := os.OpenFile(dir, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("folder -> %w", err)
	}
	defer file.Close()
	_, err = file.Write(v)
	if err != nil {
		return fmt.Errorf("folder -> %w", err)
	}
	return nil
}

func (s *Store[V]) Get(k string) ([]byte, error) {
	dir := path.Join(s.dir, k)
	file, err := os.Open(dir)
	if err != nil {
		return []byte{}, fmt.Errorf("folder -> %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, fmt.Errorf("folder -> %w", err)
	}
	return content, nil
}

func (s *Store[V]) Delete(k string) error {
	dir := path.Join(s.dir, k)
	err := os.Remove(dir)
	if err != nil {
		return fmt.Errorf("folder -> %w", err)
	}
	return nil
}

func (s *Store[V]) Clear() error {
	err := os.RemoveAll(s.dir)
	if err != nil {
		return fmt.Errorf("folder -> %w", err)
	}
	err = os.Mkdir(s.dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("folder -> %w", err)
	}
	return nil
}

func (s *Store[V]) Close() error {
	err := os.RemoveAll(s.dir)
	if err != nil {
		return fmt.Errorf("folder -> %w", err)
	}
	return nil
}
