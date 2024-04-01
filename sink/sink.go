package sink

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/wanmail/alert-fetcher/label"
	"github.com/wanmail/alert-fetcher/sink/alertmanager"
)

type SinkConfig struct {
	SinkType   string          `json:"sinkType"`
	SinkConfig json.RawMessage `json:"sinkConfig"`
}

type Sink interface {
	Send(label.Message) error
	AsyncSend(<-chan label.Message)
}

func New(cfg SinkConfig) (sink Sink, err error) {
	switch cfg.SinkType {
	case "alertmanager":
		c := alertmanager.ClientConfig{}
		if err = json.Unmarshal(cfg.SinkConfig, &c); err != nil {
			return
		}
		return alertmanager.NewClient(c)

	default:
		return nil, errors.Errorf("invalid source type %s", cfg.SinkType)
	}
}
