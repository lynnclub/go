package pool

import (
	"sync"
)

type Pool[T any] struct {
	Pool   sync.Map
	Create func(key any) *T
	Close  func(key, value any) bool
	mutex  sync.Mutex
}

func (p *Pool[T]) Get(key any) *T {
	if instance, ok := p.Pool.Load(key); ok {
		return instance.(*T)
	} else {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if instance, ok = p.Pool.Load(key); ok {
			return instance.(*T)
		}
	}

	newInstance := p.Create(key)

	p.Pool.Store(key, newInstance)
	return newInstance
}

func (p *Pool[T]) CloseAll() {
	p.Pool.Range(p.Close)
}
