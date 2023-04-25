package graph

import (
	"context"
	"errors"
	"testing"
	"time"

	"strv-template-backend-go-api/api/graphql/middleware"
	apierrors "strv-template-backend-go-api/types/errors"

	netx "go.strv.io/net"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

var (
	errTest = errors.New("test error")
)

func Test_ErrorPresenter(t *testing.T) {
	requestID := netx.NewRequestID()

	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "api-error:unknown",
			args: args{
				ctx: func() context.Context {
					ctx := netx.WithRequestID(context.Background(), requestID)
					ctx = graphql.WithOperationContext(ctx, &graphql.OperationContext{})
					return middleware.WithStartTime(ctx, time.Now())
				}(),
				err: apierrors.NewUnknownError(errTest, "").WithPublicMessage(errTest.Error()),
			},
			expectedErr: func() error {
				err := apierrors.NewUnknownError(errTest, "").WithPublicMessage(errTest.Error())
				expectedErr := gqlerror.WrapPath(nil, err)
				expectedErr.Extensions = map[string]any{
					"code":      apierrors.CodeUnknown,
					"requestId": requestID,
				}
				return expectedErr
			}(),
		},
		{
			name: "api-error:not-found",
			args: args{
				ctx: graphql.WithOperationContext(context.Background(), &graphql.OperationContext{}),
				err: apierrors.NewNotFoundError(errTest, "").WithPublicMessage(errTest.Error()),
			},
			expectedErr: func() error {
				err := apierrors.NewNotFoundError(errTest, "").WithPublicMessage(errTest.Error())
				expectedErr := gqlerror.WrapPath(nil, err)
				expectedErr.Extensions = map[string]any{
					"code": apierrors.CodeNotFound,
				}
				return expectedErr
			}(),
		},
		{
			name: "unknown-error",
			args: args{
				ctx: graphql.WithOperationContext(context.Background(), &graphql.OperationContext{}),
				err: errTest,
			},
			expectedErr: func() error {
				expectedErr := gqlerror.WrapPath(nil, errTest)
				expectedErr.Message = "internal server error"
				expectedErr.Extensions = map[string]any{
					"code": apierrors.CodeUnknown,
				}
				return expectedErr
			}(),
		},
		{
			name: "native-graphql-error",
			args: args{
				ctx: context.Background(),
				err: &gqlerror.Error{},
			},
			expectedErr: &gqlerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presenter := ErrorPresenter(zap.NewNop())
			result := presenter(tt.args.ctx, tt.args.err)
			assert.Equal(t, tt.expectedErr, result)
		})
	}
}
