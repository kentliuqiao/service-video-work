package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"
)

type TraceIDFn func(context.Context) string

type Logger struct {
	handler   slog.Handler
	traceIDFn TraceIDFn
}

func New(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, Events{})
}

func NewWithEvents(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn,
	events Events) *Logger {
	return new(w, minLevel, serviceName, traceIDFn, events)
}

func NewWithHandler(h slog.Handler) *Logger {
	return &Logger{handler: h}
}

func NewStdLogger(l *Logger, level Level) *log.Logger {
	return slog.NewLogLogger(l.handler, slog.Level(level))
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelDebug, 3, msg, args...)
}

func (l *Logger) Debugc(ctx context.Context, caller int, msg string, args ...any) {
	l.write(ctx, LevelDebug, caller, msg, args...)
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelInfo, 3, msg, args...)
}

func (l *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	l.write(ctx, LevelInfo, caller, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelWarn, 3, msg, args...)
}

func (l *Logger) Warnc(ctx context.Context, caller int, msg string, args ...any) {
	l.write(ctx, LevelWarn, caller, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.write(ctx, LevelError, 3, msg, args...)
}

func (l *Logger) Errorc(ctx context.Context, caller int, msg string, args ...any) {
	l.write(ctx, LevelError, caller, msg, args...)
}

func (l *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	sl := slog.Level(level)

	if !l.handler.Enabled(ctx, sl) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), sl, msg, pcs[0])

	if l.traceIDFn != nil {
		args = append(args, "trace_id", l.traceIDFn(ctx))
	}
	r.Add(args...)

	l.handler.Handle(ctx, r)
}

func new(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {
	// Convert the file name to just the name.ext when this key/value will
	// be logged.
	f := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		}

		return a
	}

	// Construct the slog JSON handler for use.
	handler := slog.Handler(slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{AddSource: true, Level: slog.Level(minLevel), ReplaceAttr: f}),
	)

	if events.Debug != nil || events.Info != nil || events.Warn != nil || events.Error != nil {
		handler = newLogHandler(handler, events)
	}

	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}

	handler = handler.WithAttrs(attrs)

	return &Logger{handler: handler, traceIDFn: traceIDFn}
}
