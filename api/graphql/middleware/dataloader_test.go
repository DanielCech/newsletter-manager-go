package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"strv-template-backend-go-api/database/sql"
	userdataloader "strv-template-backend-go-api/domain/user/postgres/dataloader"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDataContext struct{}

func (m mockDataContext) Ctx() context.Context {
	return context.TODO()
}

type mockDataSource struct{}

func (mockDataSource) AcquireConnCtx(_ context.Context) (sql.DataContext, error) {
	return mockDataContext{}, nil
}

func (mockDataSource) ReleaseConnCtx(_ sql.DataContext) error {
	return nil
}

func (mockDataSource) Begin(_ context.Context) (sql.DataContext, error) {
	return mockDataContext{}, nil
}

func (mockDataSource) Commit(_ sql.DataContext) error {
	return nil
}

func (mockDataSource) Rollback(_ sql.DataContext) error {
	return nil
}

func Test_DataLoader(t *testing.T) {
	type args struct {
		dataSource mockDataSource
	}
	tests := []struct {
		name               string
		args               args
		handler            http.Handler
		request            *http.Request
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "success",
			args: args{dataSource: mockDataSource{}},
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Helper()
				userLoader, ok := userdataloader.LoaderFromCtx(r.Context())
				assert.True(t, ok)
				assert.NotEmpty(t, userLoader)
				w.WriteHeader(http.StatusNoContent)
			}),
			request: func() *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/test", http.NoBody)
				require.NoError(t, err)
				return r
			}(),
			expectedStatusCode: http.StatusNoContent,
			expectedBody:       http.NoBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataLoaderMiddleware := DataLoader(tt.args.dataSource)
			w := httptest.NewRecorder()
			dataLoaderMiddleware(tt.handler).ServeHTTP(w, tt.request)
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assertResponseBody(t, tt.expectedBody, w.Body)
		})
	}
}
