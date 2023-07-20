// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	video "github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	mock "github.com/stretchr/testify/mock"
)

// IUnLikeVideoUseCase is an autogenerated mock type for the IUnLikeVideoUseCase type
type IUnLikeVideoUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, cmd
func (_m *IUnLikeVideoUseCase) Execute(ctx context.Context, cmd video.UnLikeVideoCommand) (*video.UnLikeVideoResponse, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *video.UnLikeVideoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, video.UnLikeVideoCommand) (*video.UnLikeVideoResponse, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, video.UnLikeVideoCommand) *video.UnLikeVideoResponse); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*video.UnLikeVideoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, video.UnLikeVideoCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUnLikeVideoUseCase creates a new instance of IUnLikeVideoUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUnLikeVideoUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUnLikeVideoUseCase {
	mock := &IUnLikeVideoUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}