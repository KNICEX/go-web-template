package authstorage

import (
	"encoding/gob"
	"errors"
	"go-web-template/pkg/cache"
	"go-web-template/pkg/utils"
	"time"
)

type StoreMode int

const (
	ServerStore StoreMode = iota
	JwtStore
)

var (
	ErrorKeyNotFound = errors.New("key not found")
	ErrorInvalidType = errors.New("invalid value type")
)

type Item struct {
	Key      string
	Value    any
	ExpireAt time.Time
}

func init() {
	gob.Register(Item{})
}

type Storage interface {
	Set(any, time.Duration) (string, error)
	Get(string) (*Item, error)
	Delete(string) error
	SetPrefix(string)
}

type KeyGenerator func(any, time.Duration) string

func defaultKeyGenerator(value any, duration time.Duration) string {
	return utils.RandomString(32)
}

type storage struct {
	store        cache.Driver
	mode         StoreMode
	keyGenerator KeyGenerator
	prefix       string
}

func NewStorage(driver cache.Driver, mode StoreMode, prefix string, keyGenerator ...KeyGenerator) Storage {
	if len(keyGenerator) == 0 {
		keyGenerator = append(keyGenerator, defaultKeyGenerator)
	}
	return &storage{
		store:        driver,
		mode:         mode,
		prefix:       prefix,
		keyGenerator: keyGenerator[0],
	}
}

var tenYears = time.Hour * 24 * 365 * 10

func (s *storage) Set(value any, dur time.Duration) (string, error) {
	var expiredAt time.Time
	if dur > 0 {
		expiredAt = time.Now().Add(dur)
	} else {
		expiredAt = time.Now().Add(tenYears)
	}

	key := s.keyGenerator(value, dur)
	err := s.store.Set(s.prefix+key, Item{
		Key:      key,
		Value:    value,
		ExpireAt: expiredAt,
	}, dur)
	if err != nil {
		return "", err
	}
	return key, nil

}

func (s *storage) Get(key string) (*Item, error) {
	value, ok := s.store.Get(s.prefix + key)
	if !ok {
		return nil, ErrorKeyNotFound
	}
	v, ok := value.(Item)
	if !ok {
		return nil, ErrorInvalidType
	}
	return &v, nil

}

func (s *storage) Delete(key string) error {
	return s.store.Delete(key)
}

func (s *storage) SetPrefix(prefix string) {
	s.prefix = prefix
}
