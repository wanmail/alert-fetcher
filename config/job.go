package config

import (
	"github.com/wanmail/alert-fetcher/source"
)

type JobConfig struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`

	Labels     map[string]string `json:"labels"`
	Annotation map[string]string `json:"annotation"`

	Type string `json:"type"`

	SourceConfig source.SourceConfig `json:"source"`

	Sink []string `json:"sink"`
}
