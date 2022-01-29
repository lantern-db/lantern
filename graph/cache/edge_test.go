package cache

import (
	"github.com/lantern-db/lantern/graph/model"
	"testing"
	"time"
)

func TestEdgeCache_Delete(t *testing.T) {
	c := NewEdgeCache()
	t.Run("valid_case", func(t *testing.T) {
		edge1 := model.Edge{Tail: "tail", Head: "head1", Weight: 1.0, Expiration: model.NewExpiration(5 * time.Second)}
		edge2 := model.Edge{Tail: "tail", Head: "head2", Weight: 1.0, Expiration: model.NewExpiration(5 * time.Second)}
		c.Set(edge1)
		c.Set(edge2)
		if len(c.cache["tail"]) != 2 {
			t.Errorf("mismatch length")
		}
		c.Delete("tail", "head1")
		if len(c.cache["tail"]) != 1 {
			t.Errorf("mismatch_length")
		}
		c.Delete("tail", "head2")

		_, ok := c.cache["tail"]
		if ok {
			t.Errorf("keys tail shold be deleted")
		}
	})
}

func TestEdgeCache_GetAdjacent(t *testing.T) {
	c := NewEdgeCache()
	edge1 := model.Edge{Tail: "tail", Head: "head1", Weight: 1.0, Expiration: model.NewExpiration(5 * time.Second)}
	edge2 := model.Edge{Tail: "tail", Head: "head2", Weight: 1.0, Expiration: model.NewExpiration(5 * time.Second)}
	t.Run("valid_case", func(t *testing.T) {
		c.Set(edge1)
		c.Set(edge2)
		if got, found := c.GetAdjacent("tail"); !found {
			t.Errorf("not found")
		} else {
			if len(got) != 2 {
				t.Errorf("mismatch length")
			}
		}
	})
}

func TestEdgeCache_Get(t *testing.T) {
	c := NewEdgeCache()
	edge := model.Edge{Tail: "tail", Head: "head", Weight: 1.0, Expiration: model.NewExpiration(5 * time.Second)}
	t.Run("valid_case", func(t *testing.T) {
		c.Set(edge)
		got, found := c.Get("tail", "head")
		if !found {
			t.Errorf("not found")
		}
		if got.Weight != 1.0 {
			t.Errorf("GetWeight() got = %v, want %v", got, 1.0)
		}
	})

}

func TestEdgeCache_Set(t *testing.T) {
	c := NewEdgeCache()
	edge := model.Edge{Tail: "tail", Head: "head", Weight: 1.0, Expiration: model.NewExpiration(5 * time.Second)}

	t.Run("valid_case", func(t *testing.T) {
		c.Set(edge)
	})

	t.Run("not_expired", func(t *testing.T) {
		_, found := c.Get("tail", "head")
		if !found {
			t.Errorf("not found")
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("expired", func(t *testing.T) {
		_, found := c.Get("tail", "head")
		if found {
			t.Errorf("not expired")
		}
	})
}
