package cmd

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/scys12/simple-grpc-go/config"
	"github.com/scys12/simple-grpc-go/pkg/protocol/grpc"
	"github.com/scys12/simple-grpc-go/pkg/service/v1/todo"
)

func RunServer() error {
	ctx := context.Background()
	cfg := config.NewConfig()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
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

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
