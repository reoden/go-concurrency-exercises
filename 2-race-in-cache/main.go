//////////////////////////////////////////////////////////////////////
//
// Given is some code to cache key-value pairs from a database into
// the main memory (to reduce access time). Note that golang's map are
// not entirely thread safe. Multiple readers are fine, but multiple
// writers are not. Change the code to make this thread safe.
//

package main

import (
	"container/list"
	"sync"
	"testing"
)

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should gets it's content from
	Load(string) string
}

type page struct {
	Key   string
	Value string
}

type shard struct {
	cache map[string]*list.Element
	pages list.List
	mux   sync.Mutex
}

func newShard() *shard {
	return &shard{
		cache: make(map[string]*list.Element),
	}
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	shards []*shard
	load   func(string) string
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

const SHARDING_SIZE = 32

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	shards := make([]*shard, 0, SHARDING_SIZE)
	for i := 0; i < SHARDING_SIZE; i++ {
		shards = append(shards, newShard())
	}
	return &KeyStoreCache{
		load:   load.Load,
		shards: shards,
	}
}
func (k *KeyStoreCache) PageSize() int {
	l := 0
	for i := 0; i < SHARDING_SIZE; i++ {
		l += k.shards[i].pages.Len()
	}
	return l
}

func (k *KeyStoreCache) CacheSize() int {
	l := 0
	for i := 0; i < SHARDING_SIZE; i++ {
		l += len(k.shards[i].cache)
	}
	return l
}

func (k *KeyStoreCache) getShard(key string) *shard {
	h := fnv32(key)
	return k.shards[h%uint32(len(k.shards))]
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {
	c := k.getShard(key)
	c.mux.Lock()
	defer c.mux.Unlock()
	if e, ok := c.cache[key]; ok {
		c.pages.MoveToFront(e)
		return e.Value.(page).Value
	}
	// Miss - load from database and save it in cache
	p := page{key, k.load(key)}
	// if cache is full remove the least used item
	if len(c.cache) >= CacheSize {
		end := c.pages.Back()
		// remove from map
		delete(c.cache, end.Value.(page).Key)
		// remove from list
		c.pages.Remove(end)
	}
	c.pages.PushFront(p)
	c.cache[key] = c.pages.Front()
	return p.Value
}

// Loader implements KeyStoreLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run(t *testing.T) (*KeyStoreCache, *MockDB) {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache, t)

	return cache, loader.DB
}

func main() {
	run(nil)
}
