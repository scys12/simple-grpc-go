package todo

import (
	"context"
	"fmt"

	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"github.com/scys12/simple-grpc-go/pkg/logger"
	"github.com/scys12/simple-grpc-go/pkg/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_UPDATE = "UPDATE ToDo SET `Title`=?, `Description`=?, `Reminder`=? WHERE `ID`=?"

func (s *todoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	logger.Log.Info(fmt.Sprintf("Update Todo. Todo ID: %v", req.GetTodo().Id))

	span, ctx := tracer.StartSpanFromContext(ctx, "todoservice.updatetodo")
	defer span.Finish()
	span.SetTag("id", req.Todo.Id)

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if err := req.Todo.Reminder.CheckValid(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder has invalid format-> "+err.Error())
	}

	reminder := req.Todo.Reminder.AsTime()
	res, err := c.ExecContext(
		ctx,
		QUERY_UPDATE,
		req.Todo.Title,
		req.Todo.Description,
		reminder,
		req.Todo.Id,
	)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Todo.Id))
	}

	return &v1.UpdateResponse{
		ApiV:    apiVersion,
		Updated: rows,
	}, nil
}
