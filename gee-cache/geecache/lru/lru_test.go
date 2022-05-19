package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	// 这个格式？
	lru.Add("hello", String("world"))
	// v.(String) 强制类型转换
	if v, ok := lru.Get("hello"); !ok || v.(String) != "world" {
		t.Fatalf("cache hit hello = world failed!")
	}
	if v, ok := lru.Get("hello"); ok {
		t.Log("hello : " + v.(String))
	}
	if _, ok := lru.Get("hzy"); !ok {
		t.Fatalf("cache miss hzy failed!")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	// 初始化的大小只够存放两个k-v对
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if v, ok := lru.Get(k2); !ok {
		t.Fatalf("there is no key name " + k2)
	} else {
		t.Log("key2 : " + v.(String))
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	// 能放三组k-v
	cap := 3 * len("k1"+"k2")
	t.Log(cap)

	// 在删除的时候才会调用回调
	lru := New(int64(cap), callback)

	lru.Add("k1", String("v1"))
	t.Log(lru.nBytes)
	t.Log(lru.Len())
	lru.Add("k2", String("v2"))
	t.Log(lru.nBytes)
	t.Log(lru.Len())
	lru.Add("k3", String("v3"))
	t.Log(lru.nBytes)
	t.Log(lru.Len())
	lru.Add("k4", String("v4"))
	t.Log(lru.nBytes)
	t.Log(lru.Len())
	lru.Add("k5", String("v5"))
	t.Log(lru.nBytes)
	t.Log(lru.Len())

	t.Log(keys)

	expect := []string{"k1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		//t.Log(keys)
		t.Fatalf("Call onEvicted failed, expect keys equels to %s", expect)
	}

}
