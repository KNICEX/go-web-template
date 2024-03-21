package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	rc redis.Client
}

var _ Driver = &RedisStore{}

func NewRedisStore(network, addr, user, password string, database, poolSize int) *RedisStore {
	return &RedisStore{
		rc: *redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       database,
			Username: user,
			Network:  network,
			PoolSize: poolSize,
		}),
	}
}

func (store *RedisStore) Set(key string, value any, ttl time.Duration) error {
	serialized, err := serialize(value)
	if err != nil {
		return err
	}
	return store.rc.Set(context.Background(), key, serialized, ttl).Err()
}

func (store *RedisStore) Get(key string) (any, bool) {
	val, err := store.rc.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, false
	}
	finalV, err := deserialize(val)
	return finalV, true
}

func (store *RedisStore) Gets(keys []string, prefix string) (map[string]any, []string) {
	pipe := store.rc.Pipeline()
	m := make(map[string]any)
	miss := make([]string, 0)
	var res = make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		res[i] = pipe.Get(context.Background(), prefix+key)
	}
	_, err := pipe.Exec(context.Background())
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, keys
	}
	for i, key := range keys {
		val, err := res[i].Bytes()
		if err != nil {
			miss = append(miss, key)
		} else {
			finalV, err := deserialize(val)
			if err != nil {
				miss = append(miss, key)
			} else {
				m[key] = finalV
			}
		}
	}
	return m, miss
}

func (store *RedisStore) Sets(values map[string]any, prefix string, duration time.Duration) error {
	pipe := store.rc.Pipeline()
	for key, value := range values {
		serialized, err := serialize(value)
		if err != nil {
			return err
		}
		pipe.Set(context.Background(), prefix+key, serialized, duration)
	}
	_, err := pipe.Exec(context.Background())
	return err
}

func (store *RedisStore) Delete(key string) error {
	return store.rc.Del(context.Background(), key).Err()
}

func (store *RedisStore) Deletes(keys []string, prefix string) error {
	pipe := store.rc.Pipeline()
	for _, key := range keys {
		pipe.Del(context.Background(), prefix+key)
	}
	_, err := pipe.Exec(context.Background())
	return err
}

func (store *RedisStore) Persist(_ ...string) error {
	return nil
}

func (store *RedisStore) Restore(_ ...string) error {
	return nil
}

func (store *RedisStore) Close() error {
	return store.rc.Close()
}

func (store *RedisStore) Ping() error {
	return store.rc.Ping(context.Background()).Err()
}
