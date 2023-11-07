// Package folder implements key value store on local folder
package folder

import (
	"io"
	"os"
	"path"
)

type Store[V []byte] struct {
	dir string
}

func Create(dir string) *Store[[]byte] {
	os.Mkdir(dir, os.ModePerm)
	return &Store[[]byte]{dir: dir}
}

func (s *Store[V]) Set(k string, v []byte) error {
	dir := path.Join(s.dir, k)
	file, err := os.OpenFile(dir, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(v)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store[V]) Get(k string) ([]byte, error) {
	dir := path.Join(s.dir, k)
	file, err := os.Open(dir)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func (s *Store[V]) Delete(k string) error {
	dir := path.Join(s.dir, k)
	err := os.Remove(dir)
	return err
}

func (s *Store[V]) Clear() error {
	err := os.RemoveAll(s.dir)
	os.Mkdir(s.dir, os.ModePerm)
	return err
}

func (s *Store[V]) Close() error {
	err := os.RemoveAll(s.dir)
	return err
}
