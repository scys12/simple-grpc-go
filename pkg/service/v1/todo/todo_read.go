package todo

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_READ = "SELECT `ID`, `Title`, `Description`, `Reminder` FROM ToDo WHERE `ID`=?"

func (s *todoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
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
	if err := rows.Scan(todo.Id, todo.Title, todo.Description, reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
	}
	todo.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}

	return &v1.ReadResponse{
		ApiV: apiVersion,
		Todo: &todo,
	}, nil
}
