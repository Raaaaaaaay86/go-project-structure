// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	comment "github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"

	mock "github.com/stretchr/testify/mock"
)

// IVideoCommentCreateUseCase is an autogenerated mock type for the IVideoCommentCreateUseCase type
type IVideoCommentCreateUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, cmd
func (_m *IVideoCommentCreateUseCase) Execute(ctx context.Context, cmd comment.CreateCommentCommand) (*comment.CreateCommentResponse, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *comment.CreateCommentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, comment.CreateCommentCommand) (*comment.CreateCommentResponse, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, comment.CreateCommentCommand) *comment.CreateCommentResponse); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*comment.CreateCommentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, comment.CreateCommentCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIVideoCommentCreateUseCase creates a new instance of IVideoCommentCreateUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIVideoCommentCreateUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IVideoCommentCreateUseCase {
	mock := &IVideoCommentCreateUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
