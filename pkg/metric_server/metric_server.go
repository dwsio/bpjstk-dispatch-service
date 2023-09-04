package metric_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type MetricServer struct {
	server *echo.Echo
}

func NewMetricServer(path string, serviceName string) *MetricServer {
	e := echo.New()
	e.Use(echoprometheus.NewMiddleware(serviceName))
	e.GET(path, echoprometheus.NewHandler())

	return &MetricServer{
		server: e,
	}
}

func (ms *MetricServer) Run(port string) error {
	if err := ms.server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}

func (ms *MetricServer) Shutdown(ctx context.Context) error {
	if err := ms.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
