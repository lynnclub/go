package pool

import (
	"sync"
)

type Pool[T any] struct {
	Pool          sync.Map
	mutex         sync.Mutex
	Create        func(key any) *T
	LockForCreate bool // 批量Create时建议加锁
	Close         func(key, value any) bool
}

func (p *Pool[T]) Get(key any) *T {
	if instance, ok := p.Pool.Load(key); ok {
		return instance.(*T)
	} else if p.LockForCreate {
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
