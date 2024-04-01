package source

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/wanmail/alert-fetcher/source/elasticsearch"
)

type Source interface {
	StreamSource
	BatchSource
}

// output streams one by one
type StreamSource interface {
	FetchAll(ctx context.Context, from time.Time, now time.Time) ([]map[string]interface{}, error)
}

// output batch once
type BatchSource interface {
	FetchOne(ctx context.Context, from time.Time, now time.Time) (map[string]interface{}, error)
}

type SourceConfig struct {
	SourceType   string          `json:"sourceType"`
	SourceConfig json.RawMessage `json:"sourceConfig"`
}

func New(cfg SourceConfig) (source Source, err error) {
	switch cfg.SourceType {
	case "elasticsearch":
		c := elasticsearch.Config{}
		if err = json.Unmarshal(cfg.SourceConfig, &c); err != nil {
			return
		}
		return elasticsearch.NewElasticSource(c.ClientConfig, c.QueryConfig)

	default:
		return nil, errors.Errorf("invalid source type %s", cfg.SourceType)
	}

}
