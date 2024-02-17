package handlers

import (
	"github.com/kentliuqiao/service/app/services/sales-api/v1/handlers/checkgrp"
	"github.com/kentliuqiao/service/app/services/sales-api/v1/handlers/hackgrp"
	"github.com/kentliuqiao/service/app/services/sales-api/v1/handlers/usergrp"
	v1 "github.com/kentliuqiao/service/business/web/v1"
	"github.com/kentliuqiao/service/foundation/web"
)

type Routes struct{}

func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgrp.Routes(app, cfg.Auth)
	checkgrp.Routes(app, cfg.Build, cfg.Log)
	usergrp.Routes(app, cfg.Build, cfg.Log, cfg.DB, cfg.Auth)
}
