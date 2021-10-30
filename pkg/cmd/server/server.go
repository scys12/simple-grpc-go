package cmd

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/scys12/simple-grpc-go/config"
	"github.com/scys12/simple-grpc-go/pkg/gateway"
	"github.com/scys12/simple-grpc-go/pkg/logger"
	"github.com/scys12/simple-grpc-go/pkg/protocol/grpc"
	"github.com/scys12/simple-grpc-go/pkg/protocol/rest"
	"github.com/scys12/simple-grpc-go/pkg/service/v1/todo"
)

func RunServer() error {
	ctx := context.Background()
	cfg := config.NewConfig()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initiaize log: %v", err)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	param := "parseTime=true"

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		param,
	)

	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := todo.NewTodoServiceServer(db)

	opts := gateway.Options{
		Addr: ":" + cfg.HTTPPort,
		GRPCServer: gateway.Endpoint{
			Network: "tcp",
			Addr:    fmt.Sprintf("%v:%v", cfg.GRPCHost, cfg.GRPCPort),
		},
		OpenAPIDir: cfg.OpenAPIDir,
	}

	go func() {
		_ = rest.RunServer(ctx, opts)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
