package lru

/**
分布式缓存，k-v
*/

import (
	"container/list"
	"log"
)

type Cache struct {
	maxBytes int64 // 最大允许使用的内存
	nBytes   int64 // 当前已使用的内存
	ll       *list.List
	// value 是双向链表中对应节点的指针
	cache map[string]*list.Element
	// 可选，在删除缓存时执行的回调函数
	OnEvicted func(key string, value Value)
}

// 双向链表的节点
type entry struct {
	key   string
	value Value
}

// Value是自定义的类型，有属性，使用Len计算占用多少字节
type Value interface {
	Len() int
}

// 实例化cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		// ?
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 修改当前容量，减去key和value的长度
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			// ?
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	// ele : list.Element
	if ele, ok := c.cache[key]; ok {
		// 如果在链表里查找到，就说明原先就存在，只需放到链表的最头部，并覆盖原先key的value
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry) // 这个Value不是自己造的Value
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		kv := ele.Value.(*entry)
		log.Println("db中每个的bytes：", int64(len(kv.key))+int64(kv.value.Len()))
		c.nBytes += int64(len(kv.key)) + int64(kv.value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		// 超出容量，lru删除缓存
		c.RemoveOldest()
	}
}

// Len() 获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len() // 此Len非彼Len
}
