package datatypes

import (
	"sync"
)

type MutexResource[T any] struct {
	Mu  sync.RWMutex
	Val T
}

func NewMutexResource[T any](t T) *MutexResource[T] {
	return &MutexResource[T]{
		Mu:  sync.RWMutex{},
		Val: t,
	}
}

type SafeResource[T any] interface {
	Set(T)
	Get() T
	Modify(func(T) T)
}

type safeResource[T any] struct {
	*MutexResource[T]
}

func NewSafeResource[T any](t T) SafeResource[T] {
	return newSafeResource(t)
}

func newSafeResource[T any](t T) *safeResource[T] {
	return &safeResource[T]{
		MutexResource: NewMutexResource(t),
	}
}

func (r *safeResource[T]) Set(val T) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Val = val
}

func (r *safeResource[T]) Get() T {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	val := r.Val
	return val
}

func (r *safeResource[T]) Modify(modifyFn func(T) T) {
	if modifyFn == nil {
		return
	}

	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Val = modifyFn(r.Val)
}
