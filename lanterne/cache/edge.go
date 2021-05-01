package cache

import (
	"sync"
	"time"
)

type Weight struct {
	value      float32
	expiration int64
}

type EdgeCache struct {
	ttl   time.Duration
	cache map[string]map[string]*Weight
	mu    sync.RWMutex
}

func NewEdgeCache(ttl time.Duration) EdgeCache {
	return EdgeCache{
		ttl:   ttl,
		cache: make(map[string]map[string]*Weight),
	}
}

func (c *EdgeCache) Set(tail string, head string, value float32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[tail]
	if !ok {
		c.cache[tail] = make(map[string]*Weight)
	}
	c.cache[tail][head] = &Weight{
		value:      value,
		expiration: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *EdgeCache) Delete(tail string, head string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.cache[tail]; ok {
		if _, ok := c.cache[tail][head]; ok {
			delete(c.cache[tail], head)
			if len(c.cache[tail]) == 0 {
				delete(c.cache, tail)
			}
		}
	}
}

func (c *EdgeCache) GetWeight(tail string, head string) (float32, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	weight := c.cache[tail][head]
	if time.Now().Unix() > weight.expiration {
		c.Delete(tail, head)
		return 0.0, false
	}
	return weight.value, true
}

func (c *EdgeCache) GetHeads(tail string) (map[string]float32, bool) {
	result := make(map[string]float32)
	c.mu.RLock()
	defer c.mu.RUnlock()
	for head := range c.cache[tail] {
		w, found := c.GetWeight(tail, head)
		if found {
			result[head] = w
		}
	}
	return result, len(result) != 0
}
