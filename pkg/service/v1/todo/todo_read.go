package todo

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"github.com/scys12/simple-grpc-go/pkg/logger"
	"github.com/scys12/simple-grpc-go/pkg/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const QUERY_READ = "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo WHERE `ID`=?"

func (s *todoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	logger.Log.Info(fmt.Sprintf("Read Todo. Todo ID: %v", req.GetId()))

	span, ctx := tracer.StartSpanFromContext(ctx, "todoservice.readtodo")
	defer span.Finish()
	span.SetTag("id", req.Id)

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(
		ctx,
		QUERY_READ,
		req.Id,
	)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to insert todo -> "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}

	var todo v1.Todo
	var reminder time.Time
	if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
	}
	tmp := timestamppb.New(reminder)
	if err := tmp.CheckValid(); err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}

	todo.Reminder = tmp
	return &v1.ReadResponse{
		ApiV: apiVersion,
		Todo: &todo,
	}, nil
}
