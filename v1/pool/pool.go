package pool

import (
	"sync"
)

type Pool[T any] struct {
	Pool          sync.Map
	mutex         sync.Mutex
	Create        func(key any) (*T, error)
	LockForCreate bool // 批量Create时建议加锁
	Close         func(key, value any) bool
}

func (p *Pool[T]) Get(key any) (*T, error) {
	if instance, ok := p.Pool.Load(key); ok {
		return instance.(*T), nil
	} else if p.LockForCreate {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if instance, ok = p.Pool.Load(key); ok {
			return instance.(*T), nil
		}
	}

	newInstance, err := p.Create(key)
	if err != nil {
		return nil, err
	}

	p.Pool.Store(key, newInstance)
	return newInstance, nil
}

func (p *Pool[T]) CloseAll() {
	p.Pool.Range(p.Close)
}
