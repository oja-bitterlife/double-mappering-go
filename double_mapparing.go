package doublemapparing

import (
	"sync"
	"sync/atomic"
)

type DoubleBuffer[T any] struct {
	active    atomic.Pointer[T]
	mtx       sync.Mutex
	marshal   func(*T) ([]byte, error)
	unmarshal func([]byte) (*T, error)
}

func New[T any](
	marshal func(*T) ([]byte, error),
	unmarshal func([]byte) (*T, error),
) *DoubleBuffer[T] {
	dbm := &DoubleBuffer[T]{
		marshal:   marshal,
		unmarshal: unmarshal,
	}
	dbm.active.Store(new(T)) // 初期値としてゼロ値の構造体をセット
	return dbm
}

func (dbm *DoubleBuffer[T]) clone(src *T) (*T, error) {
	b, err := dbm.marshal(src)
	if err != nil {
		return nil, err
	}
	return dbm.unmarshal(b)
}

func (dbm *DoubleBuffer[T]) Update(fn func(data *T) error) error {
	dbm.mtx.Lock()
	defer dbm.mtx.Unlock()

	cloned, err := dbm.clone(dbm.active.Load())
	if err != nil {
		return err
	}

	if err := fn(cloned); err != nil {
		return err
	}

	dbm.active.Store(cloned)
	return nil
}

func (dbm *DoubleBuffer[T]) View(fn func(data *T) error) error {
	snap, err := dbm.clone(dbm.active.Load())
	if err != nil {
		return err
	}
	return fn(snap)
}

func (dbm *DoubleBuffer[T]) Raw() *T {
	return dbm.active.Load()
}

func (dbm *DoubleBuffer[T]) Bytes() ([]byte, error) {
	return dbm.marshal(dbm.active.Load())
}

func (dbm *DoubleBuffer[T]) Restore(b []byte) error {
	newData, err := dbm.unmarshal(b)
	if err != nil {
		return err
	}

	dbm.mtx.Lock()
	defer dbm.mtx.Unlock()

	dbm.active.Store(newData)
	return nil
}
