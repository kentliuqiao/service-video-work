package usergrp

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kentliuqiao/service/business/core/user"
	"github.com/kentliuqiao/service/business/core/user/stores/userdb"
	"github.com/kentliuqiao/service/business/web/v1/auth"
	"github.com/kentliuqiao/service/foundation/logger"
	"github.com/kentliuqiao/service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	usrCore := user.NewCore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))

	hdl := New(usrCore, cfg.Auth)
	app.Handle(http.MethodPost, version, "/users", hdl.Create)
}
