package doublemapparing

import (
	"sync"
	"sync/atomic"
)

type DoubleBuffer[T any] struct {
	active    atomic.Pointer[T]
	mtx       sync.Mutex
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

func (dbm *DoubleBuffer[T]) Update(fn func(data *T) error) error {
	dbm.mtx.Lock()
	defer dbm.mtx.Unlock()

	// 現在の Active をコピー
	cloned, err := dbm.clone(dbm.active.Load())
	if err != nil {
		return err
	}

	// クローン側を更新(Rowは安全)
	if err := fn(cloned); err != nil {
		return err
	}

	// 反映
	dbm.active.Store(cloned)
	return nil
}

func (dbm *DoubleBuffer[T]) clone(src *T) (*T, error) {
	b, err := dbm.marshal(src)
	if err != nil {
		return nil, err
	}
	var dst T
	if err := dbm.unmarshal(b, &dst); err != nil {
		return nil, err
	}
	return &dst, nil
}

// View: 読み取り専用（コピーを渡すので安全、階層が深くても安心）
func (dbm *DoubleBuffer[T]) View(fn func(data *T) error) error {
	snap, err := dbm.clone(dbm.active.Load())
	if err != nil {
		return err
	}
	return fn(snap)
}

// Raw: 生のMapを返す（最速。ただし読み取り専用として扱うのがマナー）
func (dbm *DoubleBuffer[T]) Raw() *T {
	return dbm.active.Load()
}
