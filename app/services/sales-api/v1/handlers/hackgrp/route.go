package hackgrp

import (
	"net/http"

	"github.com/kentliuqiao/service/business/web/v1/auth"
	"github.com/kentliuqiao/service/business/web/v1/mid"
	"github.com/kentliuqiao/service/foundation/web"
)

func Routes(app *web.App, a *auth.Auth) {
	authenticate := mid.Authenticate(a)
	ruleAdmin := mid.Authorize(a, auth.RuleAdminOnly)

	app.Handle(http.MethodGet, "/hack", Hack)
	app.Handle(http.MethodGet, "/hackauth", Hack, authenticate, ruleAdmin)
}
