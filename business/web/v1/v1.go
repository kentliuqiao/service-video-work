package v1

import (
	"os"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/kentliuqiao/service/foundation/logger"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
}

// RouteAdder defines the behavior that sets the routes to bind for an instance of the service.
type RouteAdder interface {
	Add(mux *httptreemux.ContextMux, cfg APIMuxConfig)
}

func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	routeAdder.Add(mux, cfg)

	return mux
}
