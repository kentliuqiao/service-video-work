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
func Routes(app *web.App, build string, log *logger.Logger, db *sqlx.DB, auth *auth.Auth) {
	const version = "v1"

	usrCore := user.NewCore(log, userdb.NewStore(log, db))

	hdl := New(usrCore, auth)
	app.Handle(http.MethodPost, version, "/users", hdl.Create)
}
