package redis

import (
	"errors"
)

type MaxMin struct {
	CacheKey string //缓存键名
	Name     string //名称，CacheKey+Name=唯一标识
}

// Get 获取
func (m *MaxMin) Get() int {
	id, err := Use("").
		Get(Ctx, m.CacheKey+m.Name).
		Int()
	if err == nil {
		return id
	}

	return 0
}

// SetMax 设置最大值，大于才会覆盖
func (m *MaxMin) SetMax(newId int) error {
	key := m.CacheKey + m.Name

	cache := Use("")
	id, err := cache.Get(Ctx, key).Int()
	if err != nil && err != Nil {
		return err
	}

	if newId > id {
		return cache.Set(Ctx, key, newId, 0).Err()
	} else {
		return errors.New("小于等于最大值")
	}
}

// SetMin 设置最小值，小于才会覆盖
func (m *MaxMin) SetMin(newId int) error {
	key := m.CacheKey + m.Name

	cache := Use("")
	id, err := cache.Get(Ctx, key).Int()
	if err != nil && err != Nil {
		return err
	}

	if newId < id {
		return cache.Set(Ctx, key, newId, 0).Err()
	} else {
		return errors.New("大于等于最小值")
	}
}

// Delete 删除
func (m *MaxMin) Delete() {
	Use("").Del(Ctx, m.CacheKey+m.Name)
}
