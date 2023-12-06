package grpcerrors

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FieldViolation struct {
	Field       string
	Description string
}

func NewFieldViolationErr(msg string, fields []FieldViolation) error {
	st := status.New(codes.InvalidArgument, msg)

	br := &errdetails.BadRequest{
		FieldViolations: make([]*errdetails.BadRequest_FieldViolation, len(fields)),
	}

	for i, f := range fields {
		br.FieldViolations[i] = &errdetails.BadRequest_FieldViolation{
			Field:       f.Field,
			Description: f.Description,
		}
	}

	st, err := st.WithDetails(br)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	return st.Err()
}
