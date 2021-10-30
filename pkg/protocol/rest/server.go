package rest

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scys12/simple-grpc-go/pkg/gateway"
	"github.com/scys12/simple-grpc-go/pkg/logger"
	"github.com/scys12/simple-grpc-go/pkg/monitoring"
	"go.uber.org/zap"
)

func RunServer(ctx context.Context, options gateway.Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := gateway.DialRPC(ctx, options.GRPCServer.Addr)

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			logger.Log.Error("Failed to close a client connection to the gRPC server:", zap.String("reason", err.Error()))
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/openapiv2/", gateway.OpenAPIServer(options.OpenAPIDir))
	mux.HandleFunc("/health", gateway.CheckHealth(conn))
	mux.Handle("/prometheus", promhttp.Handler())

	gtw, err := gateway.NewGateway(ctx, conn, options.Mux)
	if err != nil {
		return err
	}

	mux.Handle("/", gtw)
	wrappedMux := monitoring.Middleware(mux)

	s := &http.Server{
		Addr:    options.Addr,
		Handler: wrappedMux,
	}

	go func() {
		<-ctx.Done()
		logger.Log.Warn("Shutting down server")
		if err := s.Shutdown(context.Background()); err != nil {
			logger.Log.Error("Failed to shutdown server:", zap.String("reason", err.Error()))
		}
	}()

	logger.Log.Info("Starting listening server at", zap.String("address", options.Addr))

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log.Error("Failed to listen and serve", zap.String("reason", err.Error()))
		return err
	}

	return nil
}
