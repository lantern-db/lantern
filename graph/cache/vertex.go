package cache

import (
	. "github.com/lantern-db/lantern/graph/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

var (
	vertexTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "vertex_total",
		Help: "Number of cached vertex",
	})
)

type VertexCache struct {
	cache map[Key]Vertex
	mu    sync.RWMutex
}

func NewVertexCache() *VertexCache {
	return &VertexCache{
		cache: make(map[Key]Vertex),
	}
}

func (c *VertexCache) Put(vertex Vertex) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[vertex.Key()]; !ok {
		vertexTotal.Inc()
	}
	c.cache[vertex.Key()] = vertex
}

func (c *VertexCache) delete(key Key) {
	if _, ok := c.cache[key]; ok {
		vertexTotal.Dec()
		delete(c.cache, key)
	}
}

func (c *VertexCache) Delete(key Key) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.delete(key)
}

func (c *VertexCache) Get(key Key) (Vertex, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.cache[key]; !ok {
		return nil, false

	} else if item.Expiration().Dead() {
		go c.delete(key)
		return nil, false

	} else {
		return item, true

	}
}

func (c *VertexCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, vertex := range c.cache {
		if vertex.Expiration().Dead() {
			c.delete(key)
		}
	}
}
