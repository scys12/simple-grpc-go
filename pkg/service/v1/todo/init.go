package todo

import (
	"context"
	"database/sql"

	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const apiVersion = "v1"

type todoServiceServer struct {
	db *sql.DB
}

func NewTodoServiceServer(db *sql.DB) v1.TodoServiceServer {
	return &todoServiceServer{
		db: db,
	}
}

func (s *todoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}
