package web

import (
	"context"
	"time"
)

type ctxKey int

const key ctxKey = 1

// Values rpresents state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

func SetValues(ctx context.Context, v *Values) context.Context {
	return context.WithValue(ctx, key, v)
}

func GetValues(ctx context.Context) *Values {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return &Values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

func GetTraceID(ctx context.Context) string {
	return GetValues(ctx).TraceID
}

func GetTime(ctx context.Context) time.Time {
	return GetValues(ctx).Now
}

func SetStatusCode(ctx context.Context, statusCode int) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return
	}

	v.StatusCode = statusCode
}
