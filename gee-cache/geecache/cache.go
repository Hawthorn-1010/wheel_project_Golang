package geecache

/*
实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu
*/

import (
	"geecache/lru"
	"log"
	"sync"
)

type cache struct {
	mu sync.Mutex
	// 使用*？
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, byteView ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		// 这里的cacheBytes还没计算吧？
		c.lru = lru.New(c.cacheBytes, nil)
	}
	log.Println(c.lru.Len())

	c.lru.Add(key, byteView)
}

func (c *cache) get(key string) (byteView ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if byteView, ok := c.lru.Get(key); ok {
		// byteView再Get中返回value，强转一下
		return byteView.(ByteView), ok
	}
	return
}
