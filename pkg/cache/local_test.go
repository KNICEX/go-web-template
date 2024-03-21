package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLocalStore_Set(t *testing.T) {
	asserts := assert.New(t)
	asserts.NoError(Set("123", "321", -1))
}

func TestLocalStore_Get(t *testing.T) {
	asserts := assert.New(t)

	// 正常情况
	{
		asserts.NoError(Set("123", "321", -1))
		value, ok := Get("123")
		asserts.True(ok)
		asserts.Equal("321", value)
	}
	// Key不存在
	{
		value, ok := Get("not_exist")
		asserts.False(ok)
		asserts.Equal(nil, value)
	}

	// 存储struct
	{
		type testStruct struct {
			key int
			v   string
		}
		test := testStruct{key: 233, v: "test"}
		asserts.NoError(Set("struct", test, -1))
		val, ok := Get("struct")
		asserts.True(ok)
		res, ok := val.(testStruct)
		asserts.True(ok)
		asserts.Equal(test, res)
	}

	// 过期
	{
		asserts.NoError(Set("expire", "expire", 1))
		time.Sleep(2 * time.Second)
		value, ok := Get("expire")
		asserts.False(ok)
		asserts.Equal(nil, value)
	}
}

func TestLocalStore_Sets(t *testing.T) {
	asserts := assert.New(t)

	err := Sets(map[string]any{"3": "3", "4": "4"}, "", -1)
	asserts.NoError(err)
	value1, _ := Get("3")
	value2, _ := Get("4")
	asserts.Equal("3", value1)
	asserts.Equal("4", value2)

}

func TestLocalStore_Gets(t *testing.T) {
	asserts := assert.New(t)
	asserts.NoError(Set("test_1", "1", -1))

	values, missed := Gets([]string{"1", "2"}, "test_")
	asserts.Equal(map[string]any{"1": "1"}, values)
	asserts.Equal([]string{"2"}, missed)
}

func TestLocalStore_Delete(t *testing.T) {
	asserts := assert.New(t)

	asserts.NoError(Set("123", "321", -1))
	_, ok := Get("123")
	asserts.True(ok)

	err := Delete("123")
	asserts.NoError(err)

	_, exist := Get("123")
	asserts.False(exist)
}

func TestLocalStore_Deletes(t *testing.T) {
	asserts := assert.New(t)

	asserts.NoError(Set("123", "321", -1))
	asserts.NoError(Set("456", "654", -1))
	asserts.NoError(Set("789", "987", -1))

	_, ok := Get("123")
	asserts.True(ok)
	_, ok = Get("456")
	asserts.True(ok)
	_, ok = Get("789")
	asserts.True(ok)

	err := Deletes([]string{"123", "456"}, "")
	asserts.NoError(err)
	_, exist := Get("123")
	asserts.False(exist)
	_, exist = Get("456")
	asserts.False(exist)
	_, exist = Get("789")
	asserts.True(exist)
}

func TestLocalStore_GarbageCollect(t *testing.T) {
	asserts := assert.New(t)
	store := NewLocalStore()
	asserts.NoError(store.Set("123", "321", time.Second*2))
	asserts.NoError(store.Set("456", "654", -1))
	asserts.NoError(store.Set("789", "987", time.Second*10))

	time.Sleep(time.Second * 3)

	store.GarbageCollect()
	_, exist := store.Get("123")
	asserts.False(exist)
	_, exist = store.Get("456")
	asserts.True(exist)
	_, exist = store.Get("789")
	asserts.True(exist)

}

func TestLocalStore_Persist(t *testing.T) {
	asserts := assert.New(t)

	asserts.NoError(Set("123", "321", -1))
	asserts.NoError(Set("456", "654", -1))
	asserts.NoError(Set("789", "987", -1))

	err := Store.Persist(DefaultCacheFile)
	asserts.NoError(err)
}

func TestLocalStore_Restore(t *testing.T) {
	asserts := assert.New(t)

	err := Store.Restore(DefaultCacheFile)
	asserts.NoError(err)

	_, ok := Get("123")
	asserts.True(ok)
	_, ok = Get("456")
	asserts.True(ok)
	_, ok = Get("789")
	asserts.True(ok)
}
