package hackgrp

import (
	"net/http"

	"github.com/kentliuqiao/service/foundation/web"
)

func Routes(app *web.App) {
	app.Handle(http.MethodGet, "/hack", Hack)
}
