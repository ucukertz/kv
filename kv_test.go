package kv

import (
	"testing"

	"github.com/ucukertz/kv/folder"
	"github.com/ucukertz/kv/maprw"
)

func TestKvInterface(t *testing.T) {
	var MaprwStore Store[string] = maprw.Create[string]()
	var MaprwBstore Bstore = maprw.Create[[]byte]()
	t.Log(MaprwStore, MaprwBstore)
	t.Log("Maprw package satisfies the interface")

	var FolderStore Bstore = folder.Create("test")
	defer FolderStore.Close()
	t.Log(FolderStore)
	t.Log("Folder package satisfies the interface")
}
