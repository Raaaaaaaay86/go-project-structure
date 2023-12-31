// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	video "github.com/raaaaaaaay86/go-project-structure/internal/context/media/video"
	mock "github.com/stretchr/testify/mock"
)

// IUpdateVideoInfoUseCase is an autogenerated mock type for the IUpdateVideoInfoUseCase type
type IUpdateVideoInfoUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: cmd
func (_m *IUpdateVideoInfoUseCase) Execute(cmd video.UpdateVideoInfoCommand) (*video.UpdateVideoInfoResponse, error) {
	ret := _m.Called(cmd)

	var r0 *video.UpdateVideoInfoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(video.UpdateVideoInfoCommand) (*video.UpdateVideoInfoResponse, error)); ok {
		return rf(cmd)
	}
	if rf, ok := ret.Get(0).(func(video.UpdateVideoInfoCommand) *video.UpdateVideoInfoResponse); ok {
		r0 = rf(cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*video.UpdateVideoInfoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(video.UpdateVideoInfoCommand) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUpdateVideoInfoUseCase creates a new instance of IUpdateVideoInfoUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUpdateVideoInfoUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUpdateVideoInfoUseCase {
	mock := &IUpdateVideoInfoUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
