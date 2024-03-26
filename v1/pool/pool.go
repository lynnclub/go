package pool

import (
	"sync"
)

type Pool[T any] struct {
	Pool   *sync.Map
	Create func(key any) T
}

func (p *Pool[T]) Get(key any) T {
	if instance, ok := p.Pool.Load(key); ok {
		return instance.(T)
	}

	newInstance := p.Create(key)

	p.Pool.Store(key, newInstance)
	return newInstance
}
