package todo

import (
	"context"
	"fmt"

	v1 "github.com/scys12/simple-grpc-go/api/proto/v1"
	"github.com/scys12/simple-grpc-go/pkg/tracer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_DELETE = "DELETE FROM ToDo WHERE `ID`=?"

func (s *todoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "todoservice.deletetodo")
	defer span.Finish()
	span.SetTag("id", req.Id)

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, QUERY_DELETE, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}

	return &v1.DeleteResponse{
		ApiV:    apiVersion,
		Deleted: rows,
	}, nil
}
