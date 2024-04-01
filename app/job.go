package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"github.com/wanmail/alert-fetcher/config"
	"github.com/wanmail/alert-fetcher/label"
	"github.com/wanmail/alert-fetcher/source"
)

type Job struct {
	config    config.JobConfig
	extractor *label.FieldExtractor
	source    source.Source
	sink      map[string]chan<- label.Message
}

func NewJob(cfg config.JobConfig) (*Job, error) {
	return &Job{
		config: cfg,
	}, nil
}

func (j *Job) Start(ctx context.Context) (err error) {
	j.extractor = label.NewFieldExtractor(j.config.Labels)

	j.source, err = source.New(j.config.SourceConfig)
	if err != nil {
		return err
	}

	j.sink = make(map[string]chan<- label.Message)
	for _, sinkName := range j.config.Sink {
		sink, ok := sinkMaps[sinkName]
		if !ok {
			slog.Error("sink not found", "name", sinkName)
			continue
		}
		j.sink[sinkName] = sink
	}

	go func() {
		from := time.Now()

		timer := time.NewTicker(time.Second * time.Duration(j.config.Duration))

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				now := time.Now()
				err := j.Fetch(ctx, from, now)
				if err != nil {
					slog.Warn("job fetch failed", "name", j.config.Name, "err", err)
				}
				from = now
			}
		}
	}()

	return nil
}

func (j *Job) ExtractMessage(ctx context.Context, data map[string]interface{}) (label.Message, error) {
	labels := j.extractor.ExtractString(data)
	annotations := label.BuildAnnotations(labels, j.config.Annotation)

	return label.Message{
		ID:          j.config.Name,
		Labels:      labels,
		Annotations: annotations,
	}, nil
}

func (j *Job) Send(ctx context.Context, msg label.Message) error {
	for name, sink := range j.sink {
		sink <- msg

		slog.InfoContext(ctx, "send data success", "name", j.config.Name, "sink", name)
	}

	return nil
}

func (j *Job) FetchBatch(ctx context.Context, from time.Time, now time.Time) (err error) {
	data, err := j.source.FetchOne(ctx, from, now)
	if err != nil {
		return
	}
	slog.InfoContext(ctx, "fetch batch data success", "name", j.config.Name)

	message, err := j.ExtractMessage(ctx, data)
	if err != nil {
		return
	}
	slog.InfoContext(ctx, "extract data success", "name", j.config.Name)

	err = j.Send(ctx, message)
	if err != nil {
		return
	}

	return nil
}

func (j *Job) FetchStream(ctx context.Context, from time.Time, now time.Time) (err error) {
	data, err := j.source.FetchAll(ctx, from, now)
	if err != nil {
		return
	}
	slog.InfoContext(ctx, "fetch stream data success", "name", j.config.Name, "count", len(data))

	for _, d := range data {
		message, err := j.ExtractMessage(ctx, d)
		if err != nil {
			return err
		}
		slog.InfoContext(ctx, "extract data success", "name", j.config.Name)

		err = j.Send(ctx, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *Job) Fetch(ctx context.Context, from time.Time, now time.Time) (err error) {
	switch j.config.Type {
	case "batch":
		return j.FetchBatch(ctx, from, now)
	case "stream":
		return j.FetchStream(ctx, from, now)
	default:
		return errors.Errorf("unknown job type %s", j.config.Type)
	}
}
