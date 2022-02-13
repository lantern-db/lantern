package adapter

import (
	"fmt"
	"github.com/lantern-db/lantern/graph/model"
	"testing"
	"time"
)

func TestNewProtoVertexOf(t *testing.T) {
	t.Run("int_case", func(t *testing.T) {
		vertex, err := NewProtoVertexOf("key", 1, 1*time.Second)
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
		vertex, err := NewProtoVertexOf("key", int32(1), 1*time.Second)
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
		vertex, err := NewProtoVertexOf("key", nil, 1*time.Second)
		if err != nil {
			t.Error(err)
		}
		value, err := vertex.IsNil()
		if err != nil {
			t.Error(err)
		}
		if value != true {
			t.Errorf("Get() got = %t, want true", value)

		}
	})

	t.Run("nil_as_proto", func(t *testing.T) {
		vertex, _ := NewProtoVertexOf("key", nil, 1*time.Second)
		pv := vertex.AsProto()
		fmt.Println(pv)
		fmt.Println(NewProtoVertex(pv))
		fmt.Println(NewProtoVertex(pv).IsNil())

		v2 := model.NewEmptyVertexOf("key", 0)
		fmt.Println(v2)
		fmt.Println(v2.AsProto())
		fmt.Println(v2.IsNil())
		fmt.Println(NewProtoVertex(v2.AsProto()).IsNil())

	})
}
