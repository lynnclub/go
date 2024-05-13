package pool

import (
	"sync"
)

type Pool[T any] struct {
	Pool   sync.Map
	Create func(key any) *T
	Close  func(key, value any) bool
}

func (p *Pool[T]) Get(key any) *T {
	if instance, ok := p.Pool.Load(key); ok {
		return instance.(*T)
	} else {
		var mutex sync.Mutex
		mutex.Lock()
		defer mutex.Unlock()
	}

	newInstance := p.Create(key)

	p.Pool.Store(key, newInstance)
	return newInstance
}

func (p *Pool[T]) CloseAll() {
	p.Pool.Range(p.Close)
}
