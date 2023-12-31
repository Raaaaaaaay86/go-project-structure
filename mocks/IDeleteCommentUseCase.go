// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	comment "github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"

	mock "github.com/stretchr/testify/mock"
)

// IDeleteCommentUseCase is an autogenerated mock type for the IDeleteCommentUseCase type
type IDeleteCommentUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, cmd
func (_m *IDeleteCommentUseCase) Execute(ctx context.Context, cmd comment.DeleteCommentCommand) (*comment.DeleteCommentResponse, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *comment.DeleteCommentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, comment.DeleteCommentCommand) (*comment.DeleteCommentResponse, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, comment.DeleteCommentCommand) *comment.DeleteCommentResponse); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*comment.DeleteCommentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, comment.DeleteCommentCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIDeleteCommentUseCase creates a new instance of IDeleteCommentUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDeleteCommentUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDeleteCommentUseCase {
	mock := &IDeleteCommentUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
