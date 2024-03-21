package cache

import (
	"bytes"
	"encoding/gob"
	"io"
	"time"
)

type Driver interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (any, bool)
	// Gets 批量取值，返回成功取值的map 与不存在的值
	Gets(keys []string, prefix string) (map[string]any, []string)
	Sets(values map[string]any, prefix string, ttl time.Duration) error
	Delete(key string) error
	Deletes(keys []string, prefix string) error
	// Persist Save in-memory cache to disk
	Persist(paths ...string) error
	// Restore Restore cache from disk
	Restore(paths ...string) error
	Close() error
	Ping() error
}

func Set(key string, value any, ttl time.Duration) error {
	return Store.Set(key, value, ttl)
}

func Get(key string) (any, bool) {
	return Store.Get(key)
}

func Gets(keys []string, prefix string) (map[string]any, []string) {
	return Store.Gets(keys, prefix)
}

func Sets(values map[string]any, prefix string, ttl time.Duration) error {
	return Store.Sets(values, prefix, ttl)
}

func Delete(key string) error {
	return Store.Delete(key)
}

func Deletes(keys []string, prefix string) error {
	return Store.Deletes(keys, prefix)
}

func Persist(path ...string) error {
	return Store.Persist(path...)
}

func Restore(path ...string) error {
	return Store.Restore(path...)
}

type persistedMap map[string]itemWithTTL

type persistedItem struct {
	Value any
}

func init() {
	gob.Register(persistedMap{})
}

type Serializer func(val any) ([]byte, error)

type Deserializer func(data []byte, val any) error

var (
	serializeFunc   = serialize
	deserializeFunc = deserialize
)

func RegisterSerializer(serializer Serializer) {

}

func RegisterDeserializer(deserializer Deserializer) {

}

func serialize(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	storeValue := persistedItem{
		Value: value,
	}
	err := enc.Encode(storeValue)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func deserialize(data []byte) (any, error) {
	var buf bytes.Buffer
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	var storeValue persistedItem
	err := dec.Decode(&storeValue)
	if err != nil {
		return nil, err
	}
	return storeValue.Value, nil
}

func deserializeReader(reader io.Reader) (any, error) {

	dec := gob.NewDecoder(reader)
	var storeValue persistedItem
	err := dec.Decode(&storeValue)
	if err != nil {
		return nil, err
	}
	return storeValue.Value, nil
}
