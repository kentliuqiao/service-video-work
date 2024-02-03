package handlers

import (
	"github.com/dimfeld/httptreemux/v5"
	"github.com/kentliuqiao/service/app/services/sales-api/v1/handlers/hackgrp"
	v1 "github.com/kentliuqiao/service/business/web/v1"
)

type Routes struct{}

func (Routes) Add(mux *httptreemux.ContextMux, cfg v1.APIMuxConfig) {
	hackgrp.Routes(mux)
}
