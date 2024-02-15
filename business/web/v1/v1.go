package v1

import (
	"os"

	"github.com/kentliuqiao/service/business/web/v1/auth"
	"github.com/kentliuqiao/service/business/web/v1/mid"
	"github.com/kentliuqiao/service/foundation/logger"
	"github.com/kentliuqiao/service/foundation/web"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	Auth     *auth.Auth
}

// RouteAdder defines the behavior that sets the routes to bind for an instance of the service.
type RouteAdder interface {
	Add(app *web.App, cfg APIMuxConfig)
}

func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(),
		// make sure to put the Panic middleware last
		// so that it can catch any panics that occur in the handler immediately and recover from them.
		mid.Panics(),
	)

	routeAdder.Add(app, cfg)

	return app
}
