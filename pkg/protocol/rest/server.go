package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/scys12/simple-grpc-go/pkg/gateway"
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
			log.Printf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/openapiv2/", gateway.OpenAPIServer(options.OpenAPIDir))
	mux.HandleFunc("/health", gateway.CheckHealth(conn))

	gtw, err := gateway.NewGateway(ctx, conn, options.Mux)
	if err != nil {
		return err
	}

	mux.Handle("/", gtw)

	s := &http.Server{
		Addr:    options.Addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown server: %v", err)
		}
	}()

	log.Printf("Starting listening server at %s", options.Addr)

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve : %v", err)
		return err
	}

	return nil
}
