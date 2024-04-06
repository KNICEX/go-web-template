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
	data map[string]*itemWithTTL
	mu   sync.RWMutex
}

const DefaultCacheFile = "cache_persist.bin"

type itemWithTTL struct {
	ExpireAt time.Time
	Value    any
}

func (i *itemWithTTL) expiredBefore(t time.Time) bool {
	return !i.ExpireAt.IsZero() && i.ExpireAt.Before(t)
}

func NewLocalStore() *LocalStore {
	res := &LocalStore{
		data: make(map[string]*itemWithTTL),
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case t := <-ticker.C:
				cnt := 0
				res.mu.Lock()
				for k, v := range res.data {
					if cnt > 1000 {
						break
					}
					if v.expiredBefore(t) {
						delete(res.data, k)
					}
					cnt++
				}
				res.mu.Unlock()
			}
		}
	}()

	return res
}

func (store *LocalStore) Set(key string, value any, ttl time.Duration) error {
	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[key] = &itemWithTTL{
		ExpireAt: expireAt,
		Value:    value,
	}
	return nil
}

func (store *LocalStore) Get(key string) (any, bool) {
	store.mu.RLock()
	item, ok := store.data[key]
	store.mu.RUnlock()
	if !ok {
		return nil, false
	}

	now := time.Now()
	if item.expiredBefore(now) {
		store.mu.Lock()
		defer store.mu.Unlock()
		// double check
		item, ok = store.data[key]
		if !ok {
			return nil, false
		}
		if item.expiredBefore(now) {
			delete(store.data, key)
			return nil, false
		}
	}
	return item.Value, true
}

func (store *LocalStore) Gets(keys []string, prefix string) (map[string]any, []string) {
	res := make(map[string]any)
	var miss []string
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, key := range keys {
		if value, ok := store.data[prefix+key]; ok {
			res[key] = value
		} else {
			miss = append(miss, key)
		}
	}
	return res, miss
}

func (store *LocalStore) Sets(values map[string]any, prefix string, ttl time.Duration) error {
	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	for k, v := range values {
		store.data[prefix+k] = &itemWithTTL{
			ExpireAt: expireAt,
			Value:    v,
		}
	}
	return nil
}

func (store *LocalStore) Delete(key string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.data, key)
	return nil
}

func (store *LocalStore) Deletes(keys []string, prefix string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	for _, key := range keys {
		delete(store.data, prefix+key)
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
	t := time.Now()
	persisted := make(persistedMap)
	store.mu.RLock()
	defer store.mu.RUnlock()
	for k, v := range store.data {
		if !v.expiredBefore(t) {
			persisted[k] = *v
		}
	}

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
	t := time.Now()
	store.mu.Lock()
	defer store.mu.Unlock()
	for k, v := range pm {
		if !v.expiredBefore(t) {
			store.data[k] = &v
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
