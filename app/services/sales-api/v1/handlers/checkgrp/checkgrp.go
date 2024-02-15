package checkgrp

import (
	"context"
	"net/http"
	"os"

	"github.com/kentliuqiao/service/foundation/logger"
	"github.com/kentliuqiao/service/foundation/web"
)

type Handlers struct {
	build string
	log   *logger.Logger
}

func New(build string, log *logger.Logger) *Handlers {
	return &Handlers{
		build: build,
		log:   log,
	}
}

// Readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (h *Handlers) Readiness(ctx context.Context, w http.ResponseWriter,
	r *http.Request) error {
	status := "ok"
	statusCode := http.StatusOK
	data := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	h.log.Info(ctx, "readiness", "status", status)

	return web.Respond(ctx, w, data, statusCode)
}

// Liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (h *Handlers) Liveness(ctx context.Context, w http.ResponseWriter,
	r *http.Request) error {

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status     string `json:"status,omitempty"`
		Build      string `json:"build,omitempty"`
		Host       string `json:"host,omitempty"`
		Name       string `json:"name,omitempty"`
		PodIP      string `json:"podIP,omitempty"`
		Node       string `json:"node,omitempty"`
		Namespace  string `json:"namespace,omitempty"`
		GOMAXPROCS string `json:"GOMAXPROCS,omitempty"`
	}{
		Status:     "up",
		Build:      h.build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: os.Getenv("GOMAXPROCS"),
	}

	h.log.Info(ctx, "liveness", "status", "OK")

	// This handler provides a free timer loop.

	return web.Respond(ctx, w, data, http.StatusOK)
}
