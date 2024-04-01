package config

import (
	"github.com/wanmail/alert-fetcher/source"
)

type JobConfig struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`

	Labels       map[string]string `json:"labels"`
	StaticLabels map[string]string `json:"staticLabels"`
	Annotations  map[string]string `json:"annotations"`

	Type string `json:"type"`

	SourceConfig source.SourceConfig `json:"source"`

	Sink []string `json:"sink"`
}
