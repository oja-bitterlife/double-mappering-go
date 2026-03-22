// ######################################################################
// スレッドセーフなトランザクションでデータを更新できるようにする仕組み。
//
//   - Update関数内でerrorを返すと、更新が行われずロールバックされます。
//   - 単純なダブルバッファだとgoroutineで死ぬ時があるので、
//     atomic.Pointerを使用してスレッドセーフなトランザクションを実現しています。
//   - せっかくシリアライズ関数を登録するんだから、
//     ファイルへの保存にも使えるようにしておきました。
//
// ######################################################################
package doublemapparing

import (
	"sync"
	"sync/atomic"
)

// **********************************************************************
// スレッドセーフなトランザクションでデータ更新できるようにする管理用構造体
type DoubleMappering[T any] struct {
	raw       atomic.Pointer[T]
	mtx       sync.Mutex
	marshal   func(*T) ([]byte, error)
	unmarshal func([]byte) (*T, error)
}

// ==================================================
// New: DoubleBufferの初期化。シリアライズ関数とデシリアライズ関数を登録する
func New[T any](
	marshal func(*T) ([]byte, error),
	unmarshal func([]byte) (*T, error),
) *DoubleMappering[T] {
	return NewFromData[T](nil, marshal, unmarshal)
}

// ==================================================
// NewFromData: 既存のデータを初期値として新しい管理構造体を作成します。
func NewFromData[T any](
	data *T,
	marshal func(*T) ([]byte, error),
	unmarshal func([]byte) (*T, error),
) *DoubleMappering[T] {
	dbm := &DoubleMappering[T]{
		marshal:   marshal,
		unmarshal: unmarshal,
	}

	// nilが渡されたらゼロ値の構造体を初期値に
	if data == nil {
		data = new(T)
	}

	dbm.raw.Store(data)
	return dbm
}

// **********************************************************************
// トランザクションの実装
// ==================================================
// clone: データのクローンを作成するためのヘルパー関数。シリアライズとデシリアライズを利用してクローンを作成する
func (dbm *DoubleMappering[T]) clone(src *T) (*T, error) {
	b, err := dbm.marshal(src)
	if err != nil {
		return nil, err
	}
	return dbm.unmarshal(b)
}

// ==================================================
// Update: データを更新するためのメソッド。 クローンを作成して更新関数に渡し、errorがなければ置き換える。
func (dbm *DoubleMappering[T]) Update(fn func(data *T) error) error {
	dbm.mtx.Lock()
	defer dbm.mtx.Unlock()

	cloned, err := dbm.clone(dbm.raw.Load())
	if err != nil {
		return err
	}

	if err := fn(cloned); err != nil {
		return err
	}

	dbm.raw.Store(cloned)
	return nil
}

// **********************************************************************
// データの取得関数
// ==================================================
// View: クローンデータを取得する。適当に扱って壊しても大丈夫なので普段はこっち
func (dbm *DoubleMappering[T]) View() (*T, error) {
	return dbm.clone(dbm.raw.Load())
}

// ==================================================
// Raw: 生データを取得する。シングルスレッド下や高速化が必要な場面で直接データにアクセスしたい場合に使用する
func (dbm *DoubleMappering[T]) Raw() *T {
	return dbm.raw.Load()
}

// **********************************************************************
// ファイルへの保存と復元のためのメソッド
// ==================================================
// Bytes: データをシリアライズしてバイト列として取得する。ファイルに保存するために使用する
func (dbm *DoubleMappering[T]) Bytes() ([]byte, error) {
	return dbm.marshal(dbm.raw.Load())
}

// ==================================================
// Restore: バイト列からデータを復元する。ファイルから読み込んだデータを復元するために使用する
func (dbm *DoubleMappering[T]) Restore(b []byte) error {
	newData, err := dbm.unmarshal(b)
	if err != nil {
		return err
	}

	dbm.mtx.Lock()
	defer dbm.mtx.Unlock()

	dbm.raw.Store(newData)
	return nil
}
