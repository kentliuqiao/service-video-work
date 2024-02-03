package v1

import (
	"net/http"
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

func APIMux(cfg APIMuxConfig) *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"alive": true}`))
	}

	mux.Handle(http.MethodGet, "/hack", h)

	return mux
}
