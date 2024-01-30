package logger

import (
	"context"
	"log/slog"
	"time"
)

type Level slog.Level

var (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// Record represents the data that is being logged.
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

// EventFn is a function to be executed when configured against a log level.
type EventFn func(ctx context.Context, r Record)

// Events contains an assignment of an event function to a log level.
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}

func toRecord(r slog.Record) Record {
	attrs := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: attrs,
	}
}
