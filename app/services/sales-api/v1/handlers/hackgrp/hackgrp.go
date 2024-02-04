package hackgrp

import (
	"context"
	"net/http"

	"github.com/kentliuqiao/service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
