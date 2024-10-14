package middleware

import (
	"testing"
	"context"
	"encore.dev"
	"encore.dev/middleware"
	"encore.dev/beta/errs"
	"github.com/stretchr/testify/assert"
)


// MockNext is a mock implementation of middleware.Next
func MockNext(req middleware.Request) middleware.Response {
	return middleware.Response{Err: nil}
}

func TestValidationMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		payload  interface{}
		expected errs.ErrCode
	}{
		{
			name: "Valid payload",
			payload: struct {
				Name string `validate:"required"`
			}{Name: "Valid Name"},
			expected: errs.OK,
		},
		{
			name: "Invalid payload",
			payload: struct {
				Name string `validate:"required"`
			}{Name: ""},
			expected: errs.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := middleware.NewRequest(context.Background(), &encore.Request{Payload: tt.payload})
			resp := ValidationMiddleware(req, MockNext)

			if tt.expected != errs.OK {
				assert.NotNil(t, resp.Err)
				assert.Equal(t, tt.expected, errs.Code(resp.Err))
			} else {
				assert.Nil(t, resp.Err)
			}
		})
	}
}