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
)

type Jfile[V any] struct {
	file *os.File
	json map[string]V
	lock *sync.RWMutex
}

func Create[V any](dir string, name string) (*Jfile[V], error) {
	os.Mkdir(dir, os.ModePerm)
	fdir := path.Join(dir, name)
	if !strings.HasSuffix(name, ".json") {
		name += ".json"
	}
	file, err := os.OpenFile(fdir, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return &Jfile[V]{}, err
	}

	j := &Jfile[V]{file: file, json: map[string]V{}, lock: &sync.RWMutex{}}
	content, err := io.ReadAll(file)
	if err != nil {
		content = []byte("{}")
		_, err = file.Write(content)
		if err != nil {
			return &Jfile[V]{}, err
		}
		return j, nil
	}
	err = json.Unmarshal(content, &j.json)
	if err != nil {
		content = []byte("{}")
		_, err = file.Write(content)
		if err != nil {
			return &Jfile[V]{}, err
		}
	}
	return j, nil
}

func (j *Jfile[V]) Set(k string, v V) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	j.json[k] = v
	content, _ := json.Marshal(j.json)
	_, err := j.file.Write(content)
	return err
}

func (j *Jfile[V]) Get(k string) (V, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	v, ok := j.json[k]
	var err error
	if !ok {
		err = fmt.Errorf("Reading key %s failed", k)
	}
	return v, err
}

func (j *Jfile[V]) Delete(k string) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	delete(j.json, k)
	content, _ := json.Marshal(j.json)
	_, err := j.file.Write(content)
	return err
}

func (j *Jfile[V]) Clear() error {
	j.lock.Lock()
	defer j.lock.Unlock()

	clear(j.json)
	content := []byte("{}")
	_, err := j.file.Write(content)
	return err
}

func (j *Jfile[V]) Close() error {
	j.lock.Lock()
	defer j.lock.Unlock()

	clear(j.json)
	return j.file.Close()
}
