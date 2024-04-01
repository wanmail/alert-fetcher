package config

import (
	"encoding/json"
	"os"

	"github.com/wanmail/alert-fetcher/logger"
	"github.com/wanmail/alert-fetcher/sink"
)

// TODO: config reload
// TODO: config format

type AppConfig struct {
	Log logger.LogConfig `json:"log"`

	Sink map[string]sink.SinkConfig `json:"sink"`

	Job []JobConfig `json:"job"`
}

func LoadConfig(path string) AppConfig {
	var cfg AppConfig
	raw, err := os.ReadFile(path)
	if err != nil {
		panic("Error occured while reading config")

	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
