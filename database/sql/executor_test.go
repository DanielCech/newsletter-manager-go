package sql

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errDefer = errors.New("defer error")
)

func Test_WithConnection(t *testing.T) {
	dctx := dbContext{Context: context.TODO()}
	mockDataSource := &mockDataSource{}
	mockDataSource.On("AcquireConnCtx").Return(dctx, nil).Once()
	mockDataSource.On("ReleaseConnCtx", dctx).Return(nil).Once()
	f := func(DataContext) error { return nil }
	err := WithConnection(context.TODO(), mockDataSource, f)
	assert.NoError(t, err)

	mockDataSource.On("AcquireConnCtx").Return(dctx, nil).Once()
	mockDataSource.On("ReleaseConnCtx", dctx).Return(errDefer).Once()
	f = func(DataContext) error { return errTest }
	err = WithConnection(context.TODO(), mockDataSource, f)
	assert.Equal(t, errors.Join(errDefer, errTest), err)

	mockDataSource.On("AcquireConnCtx").Return(dctx, errTest).Once()
	err = WithConnection(context.TODO(), mockDataSource, nil)
	assert.EqualError(t, err, "acquiring connection: test error")

	mockDataSource.AssertExpectations(t)
}

func Test_WithConnectionResult(t *testing.T) {
	dctx := dbContext{Context: context.TODO()}
	mockDataSource := &mockDataSource{}
	mockDataSource.On("AcquireConnCtx").Return(dctx, nil).Once()
	mockDataSource.On("ReleaseConnCtx", dctx).Return(nil).Once()
	userName := "Jozko"
	f := func(DataContext) (string, error) { return userName, nil }
	result, err := WithConnectionResult[string](context.TODO(), mockDataSource, f)
	assert.NoError(t, err)
	assert.Equal(t, userName, result)

	mockDataSource.On("AcquireConnCtx").Return(dctx, nil).Once()
	mockDataSource.On("ReleaseConnCtx", dctx).Return(errDefer).Once()
	f = func(DataContext) (string, error) { return "", errTest }
	result, err = WithConnectionResult[string](context.TODO(), mockDataSource, f)
	assert.Equal(t, errors.Join(errDefer, errTest), err)
	assert.Empty(t, result)

	mockDataSource.On("AcquireConnCtx").Return(dctx, errTest).Once()
	result, err = WithConnectionResult[string](context.TODO(), mockDataSource, nil)
	assert.EqualError(t, err, "acquiring connection: test error")
	assert.Empty(t, result)

	mockDataSource.AssertExpectations(t)
}

func Test_WithTransaction(t *testing.T) {
	dctx := dbContext{Context: context.TODO()}
	mockDataSource := &mockDataSource{}
	mockDataSource.On("Begin").Return(dctx, nil).Once()
	mockDataSource.On("Commit", dctx).Return(nil).Once()
	mockDataSource.On("Rollback", dctx).Return(nil).Once()
	f := func(DataContext) error { return nil }
	err := WithTransaction(context.TODO(), mockDataSource, f)
	assert.NoError(t, err)

	mockDataSource.On("Begin").Return(dctx, nil).Once()
	mockDataSource.On("Commit", dctx).Return(errTest).Once()
	mockDataSource.On("Rollback", dctx).Return(nil).Once()
	f = func(DataContext) error { return nil }
	err = WithTransaction(context.TODO(), mockDataSource, f)
	assert.EqualError(t, err, "committing transaction: test error")

	mockDataSource.On("Begin").Return(dctx, nil).Once()
	mockDataSource.On("Rollback", dctx).Return(errDefer).Once()
	f = func(DataContext) error { return errTest }
	err = WithTransaction(context.TODO(), mockDataSource, f)
	assert.Equal(t, errors.Join(errDefer, errTest), err)

	mockDataSource.On("Begin").Return(dctx, errTest).Once()
	err = WithTransaction(context.TODO(), mockDataSource, nil)
	assert.EqualError(t, err, "beginning transaction: test error")

	mockDataSource.AssertExpectations(t)
}

func Test_WithTransactionResult(t *testing.T) {
	dctx := dbContext{Context: context.TODO()}
	mockDataSource := &mockDataSource{}
	mockDataSource.On("Begin").Return(dctx, nil).Once()
	mockDataSource.On("Commit", dctx).Return(nil).Once()
	mockDataSource.On("Rollback", dctx).Return(nil).Once()
	userName := "Jozko"
	f := func(DataContext) (string, error) { return userName, nil }
	result, err := WithTransactionResult[string](context.TODO(), mockDataSource, f)
	assert.NoError(t, err)
	assert.Equal(t, userName, result)

	mockDataSource.On("Begin").Return(dctx, nil).Once()
	mockDataSource.On("Commit", dctx).Return(errTest).Once()
	mockDataSource.On("Rollback", dctx).Return(nil).Once()
	f = func(DataContext) (string, error) { return userName, nil }
	_, err = WithTransactionResult[string](context.TODO(), mockDataSource, f)
	assert.EqualError(t, err, "committing transaction: test error")

	mockDataSource.On("Begin").Return(dctx, nil).Once()
	mockDataSource.On("Rollback", dctx).Return(errDefer).Once()
	f = func(DataContext) (string, error) { return "", errTest }
	result, err = WithTransactionResult[string](context.TODO(), mockDataSource, f)
	assert.Equal(t, errors.Join(errDefer, errTest), err)
	assert.Empty(t, result)

	mockDataSource.On("Begin").Return(dctx, errTest).Once()
	result, err = WithTransactionResult[string](context.TODO(), mockDataSource, nil)
	assert.EqualError(t, err, "beginning transaction: test error")
	assert.Empty(t, result)

	mockDataSource.AssertExpectations(t)
}
