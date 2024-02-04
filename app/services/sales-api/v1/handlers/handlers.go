package handlers

import (
	"github.com/kentliuqiao/service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "github.com/kentliuqiao/service/business/web/v1"
	"github.com/kentliuqiao/service/foundation/web"
)

type Routes struct{}

func (Routes) Add(app *web.App, cfg v1.APIMuxConfig) {
	hackgrp.Routes(app)
}
