package cache

import (
	"testing"
	"time"
)

func TestEdgeCache_Delete(t *testing.T) {
	c := NewEdgeCache(5 * time.Second)
	t.Run("valid_case", func(t *testing.T) {
		c.Set("tail", "head1", 1.0)
		c.Set("tail", "head2", 1.0)
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
	c := NewEdgeCache(5 * time.Second)
	t.Run("valid_case", func(t *testing.T) {
		c.Set("tail", "head1", 1.0)
		c.Set("tail", "head2", 1.0)
		if got, found := c.GetAdjacent("tail"); !found {
			t.Errorf("not found")
		} else {
			if len(got) != 2 {
				t.Errorf("mismatch length")
			}
		}
	})
}

func TestEdgeCache_GetWeight(t *testing.T) {
	c := NewEdgeCache(5 * time.Second)
	t.Run("valid_case", func(t *testing.T) {
		c.Set("tail", "head", 1.0)
		got, found := c.GetWeight("tail", "head")
		if !found {
			t.Errorf("not found")
		}
		if got != 1.0 {
			t.Errorf("GetWeight() got = %v, want %v", got, 1.0)
		}
	})

}

func TestEdgeCache_Set(t *testing.T) {
	c := NewEdgeCache(5 * time.Second)

	t.Run("valid_case", func(t *testing.T) {
		c.Set("tail", "head", 0.0)
	})

	t.Run("not_expired", func(t *testing.T) {
		_, found := c.GetWeight("tail", "head")
		if !found {
			t.Errorf("not found")
		}
	})
	time.Sleep(6 * time.Second)
	t.Run("expired", func(t *testing.T) {
		_, found := c.GetWeight("tail", "head")
		if found {
			t.Errorf("not expired")
		}
	})
}
