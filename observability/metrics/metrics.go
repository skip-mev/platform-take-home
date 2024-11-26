package metrics

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"log"
	"net"
	"net/http"
)

func ServeMetrics(ctx context.Context, host string, port uint) error {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	server := http.Server{Handler: promhttp.Handler()}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error serving metrics: %v", err)
	}

	return nil
}
