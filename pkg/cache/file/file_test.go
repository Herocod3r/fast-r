package file

import "testing"

func TestFileSystemCache_Set(t *testing.T) {
	cacheStore := NewFileSystemCache()
	err := cacheStore.Set("hello", "world")
	if err != nil {
		t.Fail()
	}
}

func TestFileSystemCache_Get(t *testing.T) {
	cacheStore := NewFileSystemCache()
	err := cacheStore.Set("hello", "world")
	if err != nil {
		t.Fail()
	}
	val, err := cacheStore.Get("hello")
	if err != nil {
		t.Fail()
	}
	if val != "world" {
		t.Fail()
	}
}
