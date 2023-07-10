// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	video "github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	mock "github.com/stretchr/testify/mock"
)

// IVideoCreateUseCase is an autogenerated mock type for the IVideoCreateUseCase type
type IVideoCreateUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: cmd
func (_m *IVideoCreateUseCase) Execute(cmd video.CreateVideoCommand) (*video.CreateVideoResponse, error) {
	ret := _m.Called(cmd)

	var r0 *video.CreateVideoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(video.CreateVideoCommand) (*video.CreateVideoResponse, error)); ok {
		return rf(cmd)
	}
	if rf, ok := ret.Get(0).(func(video.CreateVideoCommand) *video.CreateVideoResponse); ok {
		r0 = rf(cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*video.CreateVideoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(video.CreateVideoCommand) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIVideoCreateUseCase creates a new instance of IVideoCreateUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIVideoCreateUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IVideoCreateUseCase {
	mock := &IVideoCreateUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
