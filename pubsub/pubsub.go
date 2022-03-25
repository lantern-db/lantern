package pubsub

import (
	"context"
	"log"
	"sync"
)

type Topic[T any] struct {
	mu            sync.RWMutex
	name          string
	subscriptions map[string]*Subscription[T]
}

func NewTopic[T any](name string) *Topic[T] {
	return &Topic[T]{
		name:          name,
		subscriptions: make(map[string]*Subscription[T]),
	}
}

func (t *Topic[T]) Publish(message T) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, subscription := range t.subscriptions {
		select {
		case subscription.ch <- message:
		default:
		}
	}
}

func (t *Topic[T]) NewSubscription(name string, bufferSize int) *Subscription[T] {
	t.mu.Lock()
	defer t.mu.Unlock()

	if subscription, ok := t.subscriptions[name]; !ok {
		return &Subscription[T]{
			Name:  name,
			Topic: t,
			ch:    make(chan T, bufferSize),
		}

	} else {
		return subscription

	}
}

func (t *Topic[T]) Register(subscription *Subscription[T]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subscriptions[subscription.Name]; !ok {
		t.subscriptions[subscription.Name] = subscription
	}
}

func (t *Topic[T]) Unregister(subscription *Subscription[T]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.subscriptions[subscription.Name]; ok {
		delete(t.subscriptions, subscription.Name)
	}
}

type Subscription[T any] struct {
	mu    sync.RWMutex
	Name  string
	Topic *Topic[T]
	ch    chan T
}

func (s *Subscription[T]) Register() {
	s.Topic.Register(s)
}

func (s *Subscription[T]) Unregister() {
	s.Topic.Unregister(s)
}

func (s *Subscription[T]) Subscribe(ctx context.Context, consumer func(T)) {
	var wg sync.WaitGroup
	s.Register()

	for {
		select {
		case m, ok := <-s.ch:
			if ok {
				wg.Add(1)
				go func(m T) {
					defer wg.Done()
					consumer(m)
				}(m)
			}

		case <-ctx.Done():
			s.Unregister()
			log.Printf("closing subscription: %s\n", s.Name)
			for {
				select {
				case m, ok := <-s.ch:
					if ok {
						wg.Add(1)
						go func(m T) {
							defer wg.Done()
							consumer(m)
						}(m)
					}

				default:
					wg.Wait()
					return
				}
			}
		}
	}
}
