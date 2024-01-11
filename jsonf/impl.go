// Package jsonf implements key value store on local JSON file
package jsonf

import (
	"encoding/json"
	"errors"
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

func Make[V any](dir string, name string) (*Store[V], error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if errors.Is(err, os.ErrPermission) {
		return &Store[V]{}, fmt.Errorf("%w jsonf make: %w", kv.ErrUnauthorized, err)
	} else if err != nil {
		return &Store[V]{}, fmt.Errorf("%w jsonf make: %w", kv.ErrHalt, err)
	}
	fdir := path.Join(dir, name)
	if !strings.HasSuffix(name, ".json") {
		name += ".json"
	}
	file, err := os.OpenFile(fdir, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if errors.Is(err, os.ErrPermission) {
		return &Store[V]{}, fmt.Errorf("%w jsonf make: %w", kv.ErrUnauthorized, err)
	} else if err != nil {
		return &Store[V]{}, fmt.Errorf("%w jsonf make: %w", kv.ErrHalt, err)
	}

	j := &Store[V]{file: file, m: map[string]V{}}
	content, err := io.ReadAll(file)
	if err != nil {
		content = []byte("{}")
		_, err = file.Write(content)
		if err != nil {
			return &Store[V]{}, fmt.Errorf("%w jsonf make: %w", kv.ErrHalt, err)
		}
		return j, nil
	}
	err = json.Unmarshal(content, &j.m)
	if err != nil {
		content = []byte("{}")
		_, err = file.Write(content)
		if err != nil {
			return &Store[V]{}, fmt.Errorf("%w jsonf make: %w", kv.ErrHalt, err)
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
		return fmt.Errorf("%w jsonf set %s: %w", kv.ErrHalt, k, err)
	}
	return nil
}

func (j *Store[V]) Get(k string) (V, error) {
	j.mtx.RLock()
	defer j.mtx.RUnlock()

	v, ok := j.m[k]
	var err error
	if !ok {
		err = fmt.Errorf("%w jsonf get %s: %w", kv.ErrNotFound, k, err)
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
		return fmt.Errorf("%w jsonf del %s: %w", kv.ErrHalt, k, err)
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
		return fmt.Errorf("%w jsonf cls: %w", kv.ErrHalt, err)
	}
	return nil
}

func (j *Store[V]) Purge() error {
	j.mtx.Lock()
	defer j.mtx.Unlock()

	clear(j.m)
	err := j.file.Close()
	if err != nil {
		return fmt.Errorf("%w jsonf prg: %w", kv.ErrHalt, err)
	}
	err = os.Remove(j.file.Name())
	if err != nil {
		return fmt.Errorf("%w jsonf prg: %w", kv.ErrHalt, err)
	}
	return nil
}
