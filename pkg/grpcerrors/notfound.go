package grpcerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewNotFoundErr(msg string) error {
	st := status.New(codes.NotFound, msg)
	return st.Err()
}
