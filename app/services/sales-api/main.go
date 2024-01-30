package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/kentliuqiao/service/foundation/logger"
)

var build = "develop"

func main() {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******** Send Alert ********")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return ""
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES-API", traceIDFn, events)

	ctx := context.Background()
	log.Info(ctx, "main Started")
	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		return
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// GoMAXPROCS
	log.Info(ctx, "startup", "GoMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	log.Info(ctx, "startup", "status", "initializing service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown
	log.Info(ctx, "shutdown", "status", "shutdwon start", "signal", sig)
	defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}
