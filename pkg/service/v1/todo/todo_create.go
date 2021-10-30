package todo

import (
	"context"

	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"github.com/scys12/simple-grpc-go/pkg/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_INSERT = "INSERT INTO ToDo(`Title`, `Description`, `Reminder`) VALUES(?, ?, ?)"

func (s *todoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "todoservice.createtodo")
	defer span.Finish()
	span.SetTag("todo-name", req.GetTodo().GetTitle())

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if err := req.Todo.Reminder.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder has invalid format-> "+err.Error())
	}

	reminder := req.Todo.Reminder.AsTime()
	resp, err := c.ExecContext(
		ctx,
		QUERY_INSERT,
		req.Todo.Title,
		req.Todo.Description,
		reminder,
	)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to insert todo -> "+err.Error())
	}

	id, err := resp.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to insert todo -> "+err.Error())
	}

	return &v1.CreateResponse{
		ApiV: apiVersion,
		Id:   id,
	}, nil
}
