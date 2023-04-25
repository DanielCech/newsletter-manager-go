package metrics

import (
	"context"
	"fmt"
	"sync"

	"strv-template-backend-go-api/util"

	httpx "go.strv.io/net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// once ensures that no one will override namespace and subsystem in runtime.
	once = &sync.Once{}
)

// Config contains metrics configuration.
type Config struct {
	Port      uint   `json:"port" yaml:"port" env:"METRICS_PORT" validate:"gt=0"`
	Namespace string `json:"namespace" yaml:"namespace" env:"METRICS_NAMESPACE" validate:"required"`
	Subsystem string `json:"subsystem" yaml:"subsystem" env:"METRICS_SUBSYSTEM" validate:"required"`
}

// Server contains HTTP server.
type Server struct {
	server *httpx.Server
}

// NewServer returns new instance of Server.
// Only the first call sets up namespace and subsystem.
func NewServer(cfg Config) Server {
	once.Do(func() {
		namespace = cfg.Namespace
		subsystem = cfg.Subsystem
	})

	serverConfig := httpx.ServerConfig{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: promhttp.Handler(),
		Hooks:   httpx.ServerHooks{},
		Limits:  nil,
		Logger:  util.NewServerLogger("metrics.httpx.Server"),
	}

	return Server{
		server: httpx.NewServer(&serverConfig),
	}
}

// Run starts metrics HTTP server.
func (s Server) Run(ctx context.Context) error {
	return s.server.Run(ctx)
}
