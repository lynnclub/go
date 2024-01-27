package array

import "sync"

type SliceWithMutex[T any] struct {
	mu    sync.RWMutex
	slice []T
}

func (s *SliceWithMutex[T]) Add(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slice = append(s.slice, value)
}

func (s *SliceWithMutex[T]) Get(index int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.slice[index]
}

func (s *SliceWithMutex[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.slice)
}

func (s *SliceWithMutex[T]) All() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copiedSlice := make([]T, len(s.slice))
	copy(copiedSlice, s.slice)

	return copiedSlice
}

func (s *SliceWithMutex[T]) Pop() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	copiedSlice := make([]T, len(s.slice))
	copy(copiedSlice, s.slice)

	s.slice = nil

	return copiedSlice
}
