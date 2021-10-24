package todo

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_INSERT = "INSERT INTO ToDo(`Title`, `Description`, `Reminder`) VALUES(?, ?, ?)"

func (s *todoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.Todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder has invalid format-> "+err.Error())
	}

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
