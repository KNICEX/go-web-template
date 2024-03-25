package cache

import (
	"fmt"
	"go-web-template/pkg/utils"
	"log"
	"os"
	"sync"
	"time"
)

var _ Driver = &LocalStore{}

type LocalStore struct {
	Store *sync.Map
}

const DefaultCacheFile = "cache_persist.bin"

type itemWithTTL struct {
	Expires int64
	Value   any
}

func newItem(value any, duration time.Duration) itemWithTTL {
	expires := int64(-1)
	if duration > 0 {
		expires = time.Now().Add(duration).Unix()
	}
	return itemWithTTL{
		Value:   value,
		Expires: expires,
	}
}

func getValue(item any) (any, bool) {
	var itemObj itemWithTTL
	var ok bool
	if itemObj, ok = item.(itemWithTTL); !ok {
		return item, true
	}
	if itemObj.Expires > 0 && itemObj.Expires < time.Now().Unix() {
		return nil, false
	}
	return itemObj.Value, ok
}

func (store *LocalStore) getValue(key string) (any, bool) {
	v, ok := store.Store.Load(key)
	if !ok {
		return nil, false
	}
	item, ok := v.(itemWithTTL)
	if !ok {
		return v, true
	}
	if item.Expires > 0 && item.Expires < time.Now().Unix() {
		store.Store.Delete(key)
		return nil, false
	}
	return item.Value, true
}

// GarbageCollect 回收过期的缓存
func (store *LocalStore) GarbageCollect() {
	store.Store.Range(func(key, value any) bool {
		if item, ok := value.(itemWithTTL); ok {
			if item.Expires > 0 && item.Expires < time.Now().Unix() {
				store.Store.Delete(key)
			}
		}
		return true
	})
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		Store: &sync.Map{},
	}
}

func (store *LocalStore) Set(key string, value any, ttl time.Duration) error {
	store.Store.Store(key, newItem(value, ttl))
	return nil
}

func (store *LocalStore) Get(key string) (any, bool) {
	return store.getValue(key)
}

func (store *LocalStore) Gets(keys []string, prefix string) (map[string]any, []string) {
	res := make(map[string]any)
	var miss []string
	for _, key := range keys {
		if value, ok := store.getValue(prefix + key); ok {
			res[key] = value
		} else {
			miss = append(miss, key)
		}
	}
	return res, miss
}

func (store *LocalStore) Sets(values map[string]any, prefix string, ttl time.Duration) error {
	for key, value := range values {
		store.Store.Store(prefix+key, newItem(value, ttl))
	}
	return nil
}

func (store *LocalStore) Delete(key string) error {
	store.Store.Delete(key)
	return nil
}

func (store *LocalStore) Deletes(keys []string, prefix string) error {
	for _, key := range keys {
		store.Store.Delete(prefix + key)
	}
	return nil
}

func (store *LocalStore) Persist(paths ...string) error {
	var path string
	if len(paths) > 0 {
		path = paths[0]
	} else {
		path = DefaultCacheFile
	}
	persisted := make(persistedMap)
	store.Store.Range(func(key, value any) bool {
		if item, ok := value.(itemWithTTL); ok {
			persisted[key.(string)] = item
		}
		return true
	})

	res, err := serialize(persisted)
	if err != nil {
		return fmt.Errorf("failed to serialize cache: %s", err)
	}
	err = os.WriteFile(path, res, 0644)
	return err
}

func (store *LocalStore) Restore(paths ...string) (err error) {
	var path string
	if len(paths) > 0 {
		path = paths[0]
	} else {
		path = DefaultCacheFile
	}

	if !utils.Exists(path) {
		return nil
	}
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %s", err)
	}

	defer func() {
		_ = f.Close()
		if err == nil {
			_ = os.Remove(path)
		}
	}()

	var item any
	item, err = deserializeReader(f)
	if err != nil {
		return fmt.Errorf("unknown cache file format: %s", err)
	}

	pm := item.(persistedMap)

	loaded := 0
	for k, v := range pm {
		if _, ok := getValue(v); ok {
			store.Store.Store(k, v)
			loaded++
		}
	}
	log.Println(fmt.Sprintf("Restored %d items from cache file %s", loaded, path))
	return nil
}

func (store *LocalStore) Close() error {
	return store.Persist()
}

func (store *LocalStore) Ping() error {
	return nil
}
