// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	video "github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	mock "github.com/stretchr/testify/mock"
)

// ILikeVideoUseCase is an autogenerated mock type for the ILikeVideoUseCase type
type ILikeVideoUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, cmd
func (_m *ILikeVideoUseCase) Execute(ctx context.Context, cmd video.LikeVideoCommand) (*video.LikeVideoResponse, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *video.LikeVideoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, video.LikeVideoCommand) (*video.LikeVideoResponse, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, video.LikeVideoCommand) *video.LikeVideoResponse); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*video.LikeVideoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, video.LikeVideoCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewILikeVideoUseCase creates a new instance of ILikeVideoUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewILikeVideoUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ILikeVideoUseCase {
	mock := &ILikeVideoUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
