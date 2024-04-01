package app

import (
	"log/slog"

	"github.com/wanmail/alert-fetcher/label"
	"github.com/wanmail/alert-fetcher/sink"
)

var sinkMaps = map[string]chan<- label.Message{}

func InitSink(cfgs map[string]sink.SinkConfig) {
	for name, cfg := range cfgs {
		sink, err := sink.New(cfg)
		if err != nil {
			panic(err)
		}

		ch := make(chan label.Message)
		sinkMaps[name] = ch

		go sink.AsyncSend(ch)

		slog.Info("sink init success", "name", name)
	}
}
