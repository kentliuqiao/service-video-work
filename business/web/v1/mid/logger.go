package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/kentliuqiao/service/foundation/logger"
	"github.com/kentliuqiao/service/foundation/web"
)

func Logger(log *logger.Logger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			v := web.GetValues(ctx)

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = path + "?" + r.URL.RawQuery
			}

			log.Info(ctx, "request started", "method", r.Method, "path", path,
				"remote_addr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Info(ctx, "request completed", "method", r.Method, "path", path,
				"remote_addr", r.RemoteAddr, "status", v.StatusCode, "took", time.Since(v.Now).Seconds())

			return err
		}

		return h
	}

	return m
}
