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
	config := dbm.New(
		func(v *Config) ([]byte, error) { return json.Marshal(v) },
		func(b []byte) (*Config, error) {
			var c Config
			err := json.Unmarshal(b, &c)
			return &c, err
		},
	)

	// 初期値の設定
	config.Raw().Version = 1

	view, err := config.View()
	if err != nil {
		fmt.Printf("Error viewing config: %v\n", err)
		return
	}

	config.Update(func(cfg *Config) error {
		cfg.Version = 2
		return nil
	})

	fmt.Printf("Befor Update Version: %d\n", view.Version)
	fmt.Printf("After Update Version: %d\n", config.Raw().Version)
}
