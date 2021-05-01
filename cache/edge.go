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

func (c *EdgeCache) Set(tail string, head string, value float32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[tail][head] = &Weight{
		value:      value,
		expiration: time.Now().Add(c.ttl).Unix(),
	}
}

func (c *EdgeCache) GetWeight(tail string, head string) (float32, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	weight := c.cache[tail][head]
	if time.Now().Unix() > weight.expiration {
		delete(c.cache[tail], head)
		if len(c.cache) == 0 {
			delete(c.cache, tail)
		}
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
