package authstorage

import (
	"go-web-template/pkg/cache"
	"time"
)

var defaultStorage Storage

var defaultPrefix = "auth_key."

func Init() {
	defaultStorage = NewStorage(cache.Store, ServerStore, defaultPrefix)
}

func Get(key string) (*Item, error) {
	return defaultStorage.Get(key)
}

func Set(value any, duration time.Duration) (string, error) {
	return defaultStorage.Set(value, duration)
}

func Delete(key string) error {
	return defaultStorage.Delete(key)
}
