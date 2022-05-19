package geecache

import (
	"fmt"
	"log"
	"sync"
)

/*
回调函数，命名空间
*/

type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function. ?
// 函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，
// 也能够传入实现了该接口的结构体作为参数。
type GetterFunc func(key string) ([]byte, error)

// 实现Getter的Get方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 一个 Group 可以认为是一个缓存的命名空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter can not be null!")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key can not be empty!")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache] hit!")
		return v, nil
	}
	// 没有命中缓存
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// 这个方法 ?
func (g *Group) getLocally(key string) (ByteView, error) {
	// 把回调的放入cache
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	v := ByteView{b: CloneBytes(bytes)}
	g.populateCache(key, v)
	return v, nil
}

func (g *Group) populateCache(key string, byteView ByteView) {
	g.mainCache.add(key, byteView)
}
