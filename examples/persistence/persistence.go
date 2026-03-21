package main

import (
	"encoding/json"
	"fmt"

	dbm "github.com/oja-bitterlife/double-mapparing-go"
)

// テスト用のデータ構造（ネストしたもの）
type Config struct {
	AppID   string   `json:"app_id"`
	Version int      `json:"version"`
	Meta    Metadata `json:"meta"`
}

type Metadata struct {
	Owner string `json:"owner"`
}

func main() {
	// シリアライザの定義
	src := dbm.New(
		func(c *Config) ([]byte, error) { return json.Marshal(c) },
		func(b []byte) (*Config, error) {
			var c Config
			err := json.Unmarshal(b, &c)
			return &c, err
		},
	)

	// 更新して保存する流れ
	_ = src.Update(func(c *Config) error {
		c.AppID = "Persistence Example"
		return nil
	})

	// 保存 (本来は os.WriteFile など)
	data, _ := src.Bytes()
	// debug output
	fmt.Println(string(data))

	// 復元
	dest := dbm.New(
		func(c *Config) ([]byte, error) { return json.Marshal(c) },
		func(b []byte) (*Config, error) {
			var c Config
			err := json.Unmarshal(b, &c)
			return &c, err
		},
	)
	_ = dest.Restore(data)

	fmt.Println(dest.Raw().AppID) // Output: Persistence Example
	// Output: Gopher
}
