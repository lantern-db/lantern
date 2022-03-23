package pubsub

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewTopic(t *testing.T) {
	topic := &Topic[string]{
		name:          "test",
		subscriptions: make(map[string]*Subscription[string]),
	}
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Topic[string]
	}{
		{
			name: "simple_case",
			args: args{name: "test"},
			want: topic,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTopic[string](tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscription_Subscribe(t *testing.T) {
	topic := NewTopic[string]("test")
	subscription1 := topic.NewSubscription("sub1", 10000)
	subscription2 := topic.NewSubscription("sub2", 10000)
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		subscription1.Subscribe(ctx, func(message string) {
			log.Printf("this is sub1, %s", message)
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		subscription2.Subscribe(ctx, func(message string) {
			time.Sleep(3 * time.Second)
			log.Printf("this is sub2, %s", message)
		})
	}()

	// waiting subscriptions start
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		if i > 5 {
			cancel()
		}
		time.Sleep(1 * time.Millisecond)
		topic.Publish(strconv.Itoa(i))
	}
	wg.Wait()
}
