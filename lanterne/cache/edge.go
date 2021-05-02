package cache

import (
	"sync"
	"time"
)

type itemWeight struct {
	value      float32
	expiration int64
}

type EdgeCache struct {
	ttl   time.Duration
	cache map[string]map[string]*itemWeight
	mu    sync.RWMutex
}

func NewEdgeCache(ttl time.Duration) EdgeCache {
	return EdgeCache{
		ttl:   ttl,
		cache: make(map[string]map[string]*itemWeight),
	}
}

func (c *EdgeCache) Set(tail string, head string, value float32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.cache[tail]
	if !ok {
		c.cache[tail] = make(map[string]*itemWeight)
	}
	c.cache[tail][head] = &itemWeight{
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

	if _, ok := c.cache[tail]; !ok {
		c.mu.RUnlock()
		return 0.0, false
	}

	weight, ok := c.cache[tail][head]
	c.mu.RUnlock()

	if !ok {
		return 0.0, false
	}

	if time.Now().Unix() > weight.expiration {
		c.Delete(tail, head)
		return 0.0, false
	}
	return weight.value, true
}

func (c *EdgeCache) GetAdjacent(tail string) (map[string]float32, bool) {
	result := make(map[string]float32)
	var expired []string

	c.mu.RLock()
	heads, ok := c.cache[tail]
	if !ok {
		return nil, false
	}
	for head, weight := range heads {
		if time.Now().Unix() > weight.expiration {
			expired = append(expired, head)
		} else {
			result[head] = weight.value
		}
	}
	c.mu.RUnlock()

	defer func() {
		for _, head := range expired {
			c.Delete(tail, head)
		}
	}()

	return result, len(result) != 0
}
