package app

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wanmail/alert-fetcher/config"
	"github.com/wanmail/alert-fetcher/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config file path")
}

func Start() {
	cfg := config.LoadConfig(configPath)

	logger.InitLogger(cfg.Log)

	slog.Info("config load success")

	InitSink(cfg.Sink)

	for _, jobcfg := range cfg.Job {
		job, err := NewJob(jobcfg)
		if err != nil {
			panic(err)
		}

		err = job.Start(context.Background())
		if err != nil {
			panic(err)
		}

		slog.Info("job start success", "name", jobcfg.Name)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			os.Exit(0)
		default:
			fmt.Println("signal", s)
		}
	}
}
