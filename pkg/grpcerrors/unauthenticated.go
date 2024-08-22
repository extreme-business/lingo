package grpcerrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewUnauthenticatedErr(msg string) error {
	st := status.New(codes.Unauthenticated, msg)
	return st.Err()
}
