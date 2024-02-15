package hackgrp

import (
	"net/http"

	"github.com/kentliuqiao/service/business/web/v1/auth"
	"github.com/kentliuqiao/service/business/web/v1/mid"
	"github.com/kentliuqiao/service/foundation/web"
)

func Routes(app *web.App, a *auth.Auth) {
	const version = "v1"

	authenticate := mid.Authenticate(a)
	ruleAdmin := mid.Authorize(a, auth.RuleAdminOnly)

	app.Handle(http.MethodGet, version, "/hack", Hack)
	app.Handle(http.MethodGet, version, "/hackauth", Hack, authenticate, ruleAdmin)
}
