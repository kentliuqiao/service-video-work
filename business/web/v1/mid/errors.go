package mid

import (
	"context"
	"net/http"

	"github.com/kentliuqiao/service/business/web/v1/auth"
	"github.com/kentliuqiao/service/business/web/v1/response"
	"github.com/kentliuqiao/service/foundation/logger"
	"github.com/kentliuqiao/service/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Error(ctx, "message", "err", err)

				var ed response.ErrorDocument
				var status int

				switch {
				case response.IsError(err):
					// trusted error.
					reqErr := response.GetError(err)
					ed = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					ed = response.ErrorDocument{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					// untrusted error.
					ed = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, ed, status); err != nil {
					return err
				}

				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
