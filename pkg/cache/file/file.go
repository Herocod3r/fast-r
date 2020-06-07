package file

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/herocod3r/fast-r/pkg/cache"
)

type fileSystemCache struct {
	filePath string
}

func NewFileSystemCache() cache.Client {
	dir, er := os.UserCacheDir()
	if er != nil {
		return &fileSystemCache{}
	}
	path := filepath.Join(dir, "fast_r_cache.txt")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, er := os.Create(path)
		if er != nil {
			return &fileSystemCache{}
		}
		_ = file.Close()
	}

	return &fileSystemCache{filePath: path}
}

func (f *fileSystemCache) getDictStore() (map[string]string, error) {
	if len(f.filePath) < 1 {
		return nil, cache.StoreNotActiveErr
	}
	dict := make(map[string]string)
	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, cache.ValueNotFoundErr
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)

	_ = json.Unmarshal(data, &dict)
	return dict, nil
}

func (f *fileSystemCache) Get(key string) (string, error) {

	store, er := f.getDictStore()
	if er != nil {
		return "", er
	}

	if value, ok := store[key]; ok {
		return value, nil
	}
	return "", cache.ValueNotFoundErr
}

func (f fileSystemCache) Set(key, value string) error {
	store, er := f.getDictStore()
	if er != nil {
		return er
	}

	store[key] = value
	file, err := os.OpenFile(f.filePath, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return cache.ValueNotFoundErr
	}
	defer file.Close()
	data, _ := json.Marshal(store)
	_, err = file.WriteString(string(data))
	return err
}
