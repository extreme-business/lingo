package grpcerrors_test

import (
	"testing"

	"github.com/dwethmar/lingo/pkg/grpcerrors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewFieldViolationErr(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		err := grpcerrors.NewFieldViolationErr("msg", []grpcerrors.FieldViolation{
			{
				Field:       "field",
				Description: "description",
			},
		})
		if err == nil {
			t.Errorf("err should not be nil")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Errorf("expected status to be ok")
		}

		if st.Code() != codes.InvalidArgument {
			t.Errorf("expected code to be InvalidArgument")
		}

		if len(st.Details()) != 1 {
			t.Errorf("expected 1 detail")
		}

		for _, detail := range st.Details() {
			switch v := detail.(type) {
			case *errdetails.BadRequest:
				if len(v.FieldViolations) != 1 {
					t.Errorf("expected 1 field violation")
				}
				if v.FieldViolations[0].Field != "field" {
					t.Errorf("expected field to be field")
				}
				if v.FieldViolations[0].Description != "description" {
					t.Errorf("expected description to be description")
				}
			default:
				t.Errorf("unexpected detail type: %T", v)
			}
		}
	})
}
