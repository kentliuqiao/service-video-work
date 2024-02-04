package hackgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/kentliuqiao/service/business/web/v1/response"
	"github.com/kentliuqiao/service/foundation/web"
)

func Hack(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100) % 2; n == 0 {
		return response.NewError(errors.New("trusted error"), http.StatusTeapot)
	}
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
