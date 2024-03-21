package cache

import (
	"encoding/gob"

	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func initRedis() Driver {
	d := NewRedisStore(
		"tcp",
		"localhost:6379",
		"",
		"",
		0,
		0,
	)
	err := d.Ping()
	if err != nil {
		panic(err)
	}
	return d
}

func TestRedisStore_Set(t *testing.T) {
	store := initRedis()
	defer store.Close()
	asserts := assert.New(t)
	asserts.NoError(store.Set("test", "test val", -1))

	// 带有TTL
	{
		asserts.NoError(store.Set("test", "test val", time.Second))

	}
	// 未注册结构体
	{
		type testStruct struct {
			Key int
			V   string
		}
		test := testStruct{Key: 233, V: "test"}
		asserts.Error(store.Set("struct", test, -1))

	}
	// 注册结构体
	{
		type testStruct struct {
			Key int
			V   string
		}
		gob.Register(testStruct{})
		test := testStruct{Key: 233, V: "test"}
		asserts.NoError(store.Set("struct", test, -1))
	}

}

func TestRedisStore_Get(t *testing.T) {
	store := initRedis()
	defer store.Close()
	asserts := assert.New(t)

	// 正常情况
	{
		asserts.NoError(store.Set("123", "321", -1))
		value, ok := store.Get("123")
		asserts.True(ok)
		asserts.Equal("321", value)
	}
	// Key不存在
	{
		value, ok := store.Get("not_exist")
		asserts.False(ok)
		asserts.Equal(nil, value)
	}

	// 存储struct
	{
		type testStruct struct {
			Key int
			V   string
		}
		gob.Register(testStruct{})
		test := testStruct{Key: 233, V: "test"}
		asserts.NoError(store.Set("struct", test, -1))
		val, ok := store.Get("struct")
		asserts.True(ok)
		res, ok := val.(testStruct)
		asserts.True(ok)
		asserts.Equal(test, res)
	}

	// 过期
	{
		asserts.NoError(store.Set("expire", "expire", 1))
		time.Sleep(2 * time.Second)
		value, ok := store.Get("expire")
		asserts.False(ok)
		asserts.Equal(nil, value)
	}
}

func TestRedisStore_Sets(t *testing.T) {
	store := initRedis()
	defer store.Close()
	asserts := assert.New(t)

	err := store.Sets(map[string]any{"3": "3", "4": "4"}, "", -1)
	asserts.NoError(err)
	value1, _ := store.Get("3")
	value2, _ := store.Get("4")
	asserts.Equal("3", value1)
	asserts.Equal("4", value2)
}

func TestRedisStore_Gets(t *testing.T) {
	store := initRedis()
	defer store.Close()
	asserts := assert.New(t)

	asserts.NoError(store.Set("test_1", "1", -1))
	asserts.NoError(store.Set("test_2", "2", -1))
	asserts.NoError(store.Set("test_3", "3", -1))
	asserts.NoError(store.Set("test_4", "4", -1))
	// 全部命中
	{
		res, missed := store.Gets([]string{"1", "2", "3", "4"}, "test_")
		asserts.Equal(map[string]any{"1": "1", "2": "2", "3": "3", "4": "4"}, res)
		asserts.Equal([]string{}, missed)
	}

	// 部分命中
	{
		res, missed := store.Gets([]string{"1", "2", "3", "5"}, "test_")
		asserts.Equal(map[string]any{"1": "1", "2": "2", "3": "3"}, res)
		asserts.Equal([]string{"5"}, missed)
	}
}

func TestRedisStore_Delete(t *testing.T) {
	store := initRedis()
	defer store.Close()
	asserts := assert.New(t)

	asserts.NoError(store.Set("test_1", "1", -1))
	asserts.NoError(store.Set("test_2", "2", -1))
	asserts.NoError(store.Set("test_3", "3", -1))
	asserts.NoError(store.Set("test_4", "4", -1))

	asserts.NoError(store.Delete("test_1"))
	asserts.NoError(store.Delete("test_4"))

	_, ok := store.Get("test_1")
	asserts.False(ok)
	_, ok = store.Get("test_2")
	asserts.True(ok)
	_, ok = store.Get("test_3")
	asserts.True(ok)
	_, ok = store.Get("test_4")
	asserts.False(ok)
}

func TestRedisStore_Deletes(t *testing.T) {
	store := initRedis()
	defer store.Close()
	asserts := assert.New(t)

	asserts.NoError(store.Set("test_1", "1", -1))
	asserts.NoError(store.Set("test_2", "2", -1))
	asserts.NoError(store.Set("test_3", "3", -1))
	asserts.NoError(store.Set("test_4", "4", -1))

	asserts.NoError(store.Deletes([]string{"test_1", "test_4", "test_100"}, ""))
	_, ok := store.Get("test_1")
	asserts.False(ok)
	_, ok = store.Get("test_2")
	asserts.True(ok)
	_, ok = store.Get("test_3")
	asserts.True(ok)
	_, ok = store.Get("test_4")
	asserts.False(ok)
}
