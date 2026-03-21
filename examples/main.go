package main

import (
	"encoding/json"
	"fmt"

	doublemapparing "github.com/oja-bitterlife/double-mapparing-go"
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
	db := doublemapparing.New(
		func(v *Config) ([]byte, error) { return json.Marshal(v) },
		func(b []byte) (*Config, error) {
			var c Config
			err := json.Unmarshal(b, &c)
			return &c, err
		},
	)

	db.Update(func(cfg *Config) error {
		cfg.Version = 2
		return nil
	})

	fmt.Printf("Current Version: %d\n", db.Raw().Version)
}
