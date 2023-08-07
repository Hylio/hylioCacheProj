package main

import (
	"fmt"
	"log"
	"sync"
)

/*
hyliocache/
    |--lru/
        |--lru.go    // lru 缓存淘汰策略
    |--byteview.go   // 缓存值的抽象与封装
    |--cache.go      // 并发控制
    |--hyliocache.go // 负责与外部交互，控制缓存存储和获取的主流程
*/

// Getter 实现从数据源获取数据的能力
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 通过实现Get方法 使得任意函数只要通过GetterFunc的转换 就能实现Getter接口
type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

// Group 定义一块缓存空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 新建一块缓存空间
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("no getter")
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
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("hyliocache hit!")
		return v, nil
	}
	// 缓存未命中
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	// 返回一个深拷贝 而不是源数据
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache 把最近访问过的 没有在缓存中的数据 保存在缓存中
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
