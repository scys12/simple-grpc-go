package todo

import (
	"context"
	"time"

	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"github.com/scys12/simple-grpc-go/pkg/logger"
	"github.com/scys12/simple-grpc-go/pkg/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *todoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	logger.Log.Info("Read Many Todos.")

	span, ctx := tracer.StartSpanFromContext(ctx, "todoservice.readalltodos")
	defer span.Finish()

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	list := []*v1.Todo{}
	for rows.Next() {
		td := new(v1.Todo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
		}
		tmp := timestamppb.New(reminder)
		if err := tmp.CheckValid(); err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}

		td.Reminder = tmp
		list = append(list, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
	}

	return &v1.ReadAllResponse{
		ApiV:  apiVersion,
		Todos: list,
	}, nil
}
