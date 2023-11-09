// Package jsonf implements key value store on local JSON file
package jsonf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/ucukertz/kv"
)

type Store[V any] struct {
	file *os.File
	m    map[string]V
	mtx  sync.RWMutex
}

var _ kv.Store[any] = (*Store[any])(nil)
var _ kv.Bstore = (*Store[[]byte])(nil)

func Create[V any](dir string, name string) (*Store[V], error) {
	os.MkdirAll(dir, os.ModePerm)
	fdir := path.Join(dir, name)
	if !strings.HasSuffix(name, ".json") {
		name += ".json"
	}
	file, err := os.OpenFile(fdir, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return &Store[V]{}, fmt.Errorf("jsonf -> %w", err)
	}

	j := &Store[V]{file: file, m: map[string]V{}}
	content, err := io.ReadAll(file)
	if err != nil {
		content = []byte("{}")
		_, err = file.Write(content)
		if err != nil {
			return &Store[V]{}, fmt.Errorf("jsonf -> %w", err)
		}
		return j, nil
	}
	err = json.Unmarshal(content, &j.m)
	if err != nil {
		content = []byte("{}")
		_, err = file.Write(content)
		if err != nil {
			return &Store[V]{}, fmt.Errorf("jsonf -> %w", err)
		}
	}
	return j, nil
}

func (j *Store[V]) Set(k string, v V) error {
	j.mtx.Lock()
	defer j.mtx.Unlock()

	j.m[k] = v
	content, _ := json.Marshal(j.m)
	_, err := j.file.Write(content)
	if err != nil {
		return fmt.Errorf("jsonf -> %w", err)
	}
	return nil
}

func (j *Store[V]) Get(k string) (V, error) {
	j.mtx.RLock()
	defer j.mtx.RUnlock()

	v, ok := j.m[k]
	var err error
	if !ok {
		err = fmt.Errorf("jsonf -> Reading key %s failed", k)
	}
	return v, err
}

func (j *Store[V]) Delete(k string) error {
	j.mtx.Lock()
	defer j.mtx.Unlock()

	delete(j.m, k)
	content, _ := json.Marshal(j.m)
	_, err := j.file.Write(content)
	if err != nil {
		return fmt.Errorf("jsonf -> %w", err)
	}
	return nil
}

func (j *Store[V]) Clear() error {
	j.mtx.Lock()
	defer j.mtx.Unlock()

	clear(j.m)
	content := []byte("{}")
	_, err := j.file.Write(content)
	if err != nil {
		return fmt.Errorf("jsonf -> %w", err)
	}
	return nil
}

func (j *Store[V]) Close() error {
	j.mtx.Lock()
	defer j.mtx.Unlock()

	clear(j.m)
	err := j.file.Close()
	if err != nil {
		return fmt.Errorf("jsonf -> %w", err)
	}
	return nil
}
