package adapter

import (
	"github.com/lantern-db/lantern/graph/model"
	"testing"
	"time"
)

func TestNewProtoVertexOf(t *testing.T) {
	t.Run("int_case", func(t *testing.T) {
		vertex, err := NewProtoVertexOf("key", 1, model.NewExpiration(1*time.Second))
		expected := 1
		if err != nil {
			t.Error(err)
		}
		value, err := vertex.IntValue()
		if err != nil {
			t.Error(err)
		}
		if value != expected {
			t.Errorf("Get() got = %c, want %c", value, expected)
		}
	})

	t.Run("int32_case", func(t *testing.T) {
		vertex, err := NewProtoVertexOf("key", int32(1), model.NewExpiration(1*time.Second))
		expected := 1
		if err != nil {
			t.Error(err)
		}
		value, err := vertex.IntValue()
		if err != nil {
			t.Error(err)
		}
		if value != expected {
			t.Errorf("Get() got = %c, want %c", value, expected)
		}
	})

	t.Run("nil_case", func(t *testing.T) {
		vertex, err := NewProtoVertexOf("key", nil, model.NewExpiration(1*time.Second))
		if err != nil {
			t.Error(err)
		}
		value, err := vertex.NilValue()
		if err != nil {
			t.Error(err)
		}
		if value != nil {
			t.Errorf("Get() got = %c, want nil", value)

		}
	})
}
