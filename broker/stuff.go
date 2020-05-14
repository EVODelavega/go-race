package broker

import (
	"context"
	"log"
	"sync"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/sub_mock.go -package mocks my.pkg/race/broker Sub
type Sub interface {
	C() chan<- interface{}
	Done() <-chan struct{}
}

type Broker struct {
	mu   sync.Mutex
	ctx  context.Context
	subs map[int]Sub
	keys []int
}

func New(ctx context.Context) *Broker {
	return &Broker{
		ctx:  ctx,
		subs: map[int]Sub{},
		keys: []int{},
	}
}

func (b *Broker) Send(v interface{}) {
	b.mu.Lock()
	go func() {
		rm := make([]int, 0, len(b.subs))
		defer func() {
			if len(rm) > 0 {
				b.unsub(rm...)
			}
			b.mu.Unlock()
		}()
		for k, s := range b.subs {
			select {
			case <-b.ctx.Done():
				return
			case <-s.Done():
				rm = append(rm, k)
			case s.C() <- v:
				continue
			default:
				log.Printf("Skipped sub %d", k)
			}
		}
	}()
}

func (b *Broker) Subscribe(s Sub) int {
	b.mu.Lock()
	k := b.key()
	b.subs[k] = s
	b.mu.Unlock()
	return k
}

func (b *Broker) Unsubscribe(k int) {
	b.mu.Lock()
	b.unsub(k)
	b.mu.Unlock()
}

func (b *Broker) key() int {
	if len(b.keys) > 0 {
		k := b.keys[0]
		b.keys = b.keys[1:]
		return k
	}
	return len(b.subs) + 1
}

func (b *Broker) unsub(keys ...int) {
	for _, k := range keys {
		if _, ok := b.subs[k]; !ok {
			return
		}
		delete(b.subs, k)
		b.keys = append(b.keys, k)
	}
}
