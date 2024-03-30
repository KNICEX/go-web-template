package authstorage

import (
	"encoding/gob"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-web-template/pkg/cache"
	"testing"
)

func TestNewServerStorageLocal(t *testing.T) {
	s := NewStorage(cache.NewLocalStore(), ServerStore, "")
	asserts := assert.New(t)
	type testStruct1 struct {
		Name string
		Age  int
	}
	gob.Register(testStruct1{})
	v1 := testStruct1{
		Name: "test",
		Age:  18,
	}
	key, err := s.Set(v1, 0)
	asserts.NoError(err)
	t.Log(fmt.Sprintf("set key: %s, v: %v", key, v1))

	item, err := s.Get(key)
	asserts.NoError(err)
	asserts.Equal(v1.Name, item.Value.(testStruct1).Name)
	asserts.Equal(v1.Age, item.Value.(testStruct1).Age)
}

func TestNewServerStorageRedis(t *testing.T) {
	s := NewStorage(cache.NewRedisStore(
		"tcp",
		"localhost:6379",
		"",
		"",
		0,
		10,
	), ServerStore, "")
	asserts := assert.New(t)
	type testStruct1 struct {
		Name string
		Age  int
	}
	gob.Register(testStruct1{})
	v1 := testStruct1{
		Name: "test",
		Age:  18,
	}
	key, err := s.Set(v1, 0)
	asserts.NoError(err)
	t.Log(fmt.Sprintf("set key: %s, v: %v", key, v1))

	item, err := s.Get(key)
	asserts.NoError(err)
	asserts.Equal(v1.Name, item.Value.(testStruct1).Name)

}
