package doublemappering

import (
	"encoding/json"
	"sync"
	"testing"

	dbm "github.com/oja-bitterlife/double-mappering-go"
)

type TestData struct {
	Counter int `json:"counter"`
}

func TestDoubleBuffer_Race(t *testing.T) {
	m := func(v *TestData) ([]byte, error) { return json.Marshal(v) }
	u := func(b []byte) (*TestData, error) {
		var d TestData
		err := json.Unmarshal(b, &d)
		return &d, err
	}
	testData := dbm.New[TestData](m, u)

	var wg sync.WaitGroup
	workers := 100
	iterations := 100

	// Writer: 複数のゴルーチンから同時に Update
	for range workers {
		wg.Go(func() {
			for range iterations {
				_ = testData.Update(func(d *TestData) error {
					d.Counter++
					return nil
				})
			}
		})
	}

	// Reader: 更新中にひたすら Raw で読み取り
	stopReader := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopReader:
				return
			default:
				_ = testData.Raw().Counter
			}
		}
	}()

	wg.Wait()
	close(stopReader)

	// 最終結果の整合性チェック
	final := testData.Raw().Counter
	expected := workers * iterations
	if final != expected {
		t.Errorf("expected %d, got %d", expected, final)
	}
}
