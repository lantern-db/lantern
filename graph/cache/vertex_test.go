package cache

import (
	"github.com/golang/mock/gomock"
	. "github.com/lantern-db/lantern/graph/model"
	"testing"
	"time"
)

func TestVertexCache_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := Vertex{Key: "key", Value: "value", Expiration: NewExpiration(3 * time.Second)}

	c := NewVertexCache()
	t.Run("valid_case", func(t *testing.T) {
		c.Set(v)
		c.Delete(v.Key)
		got, found := c.Get(v.Key)
		if found {
			t.Errorf("Get() got = %v, want %v", got, nil)
		}
	})
}

func TestVertexCache_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := Vertex{Key: "key", Value: "value", Expiration: NewExpiration(3 * time.Second)}

	c := NewVertexCache()
	t.Run("valid_case", func(t *testing.T) {
		c.Set(v)
		got, found := c.Get(v.Key)
		if !found {
			t.Errorf("Get() got = %c, want %c", got.Value, v.Value)
		}
	})

	t.Run("time_out", func(t *testing.T) {
		time.Sleep(4 * time.Second)
		got, found := c.Get(v.Key)
		if found {
			t.Errorf("Get() got = %c, want %v", got.Value, nil)
		}
	})
}

func TestVertexCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := Vertex{Key: "key", Value: "value", Expiration: NewExpiration(3 * time.Second)}

	c := NewVertexCache()
	t.Run("valid_case", func(t *testing.T) {
		c.Set(v)
		got, found := c.Get("key")
		if !found {
			t.Errorf("Get() got = %c, want %c", got.Value, v.Value)
		}
	})
}
