package grpcerrors_test

import (
	"testing"

	"github.com/dwethmar/lingo/pkg/grpcerrors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type expectedFieldViolation struct {
	Field       string
	Description string
}

func TestNewFieldViolationErr(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		expectedViolations := []expectedFieldViolation{
			{Field: "field", Description: "description"},
		}

		err := grpcerrors.NewFieldViolationErr("msg", []grpcerrors.FieldViolation{{
			Field:       "field",
			Description: "description",
		}})

		if err == nil {
			t.Errorf("err should not be nil")
		}

		assertStatus(t, err, expectedViolations)
	})
}

func assertStatus(t *testing.T, err error, expectedViolations []expectedFieldViolation) {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok {
		t.Errorf("expected status to be ok")
		return
	}

	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected code to be InvalidArgument")
	}

	if len(st.Details()) != len(expectedViolations) {
		t.Errorf("expected %d details, got %d", len(expectedViolations), len(st.Details()))
	}

	for i, detail := range st.Details() {
		assertDetail(t, detail, expectedViolations[i])
	}
}

func assertDetail(t *testing.T, detail interface{}, expected expectedFieldViolation) {
	t.Helper()
	switch v := detail.(type) {
	case *errdetails.BadRequest:
		if len(v.GetFieldViolations()) != 1 {
			t.Errorf("expected 1 field violation")
		}
		if v.GetFieldViolations()[0].GetField() != expected.Field {
			t.Errorf("expected field to be %s, got %s", expected.Field, v.GetFieldViolations()[0].GetField())
		}
		if v.GetFieldViolations()[0].GetDescription() != expected.Description {
			t.Errorf("expected description to be %s, got %s", expected.Description, v.GetFieldViolations()[0].GetDescription())
		}
	default:
		t.Errorf("unexpected detail type: %T", v)
	}
}
