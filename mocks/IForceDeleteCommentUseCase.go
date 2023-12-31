// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	comment "github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"

	mock "github.com/stretchr/testify/mock"
)

// IForceDeleteCommentUseCase is an autogenerated mock type for the IForceDeleteCommentUseCase type
type IForceDeleteCommentUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, cmd
func (_m *IForceDeleteCommentUseCase) Execute(ctx context.Context, cmd comment.ForceDeleteCommentCommand) (*comment.ForceDeleteCommentResponse, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *comment.ForceDeleteCommentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, comment.ForceDeleteCommentCommand) (*comment.ForceDeleteCommentResponse, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, comment.ForceDeleteCommentCommand) *comment.ForceDeleteCommentResponse); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*comment.ForceDeleteCommentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, comment.ForceDeleteCommentCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIForceDeleteCommentUseCase creates a new instance of IForceDeleteCommentUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIForceDeleteCommentUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IForceDeleteCommentUseCase {
	mock := &IForceDeleteCommentUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
