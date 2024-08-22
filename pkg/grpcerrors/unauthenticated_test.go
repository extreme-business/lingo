package grpcerrors_test

import (
	"testing"

	"github.com/extreme-business/lingo/pkg/grpcerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewUnauthenticatedErr(t *testing.T) {
	t.Run("error should match expected", func(t *testing.T) {
		msg := "invalid username or password"

		err := grpcerrors.NewUnauthenticatedErr(msg)
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected a gRPC status error, got %v", err)
		}

		// Check the status code
		if st.Code() != codes.Unauthenticated {
			t.Errorf("expected code %v, got %v", codes.Unauthenticated, st.Code())
		}

		// Check the status message
		if st.Message() != msg {
			t.Errorf("expected message %q, got %q", msg, st.Message())
		}
	})
}
