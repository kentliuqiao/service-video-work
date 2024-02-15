package checkgrp

import (
	"net/http"

	"github.com/kentliuqiao/service/foundation/logger"
	"github.com/kentliuqiao/service/foundation/web"
)

func Routes(app *web.App, build string, log *logger.Logger) {
	const version = "v1"

	hdl := New(build, log)
	app.Handle(http.MethodGet, version, "/readiness", hdl.Readiness)
	app.Handle(http.MethodGet, version, "/liveness", hdl.Liveness)
}
