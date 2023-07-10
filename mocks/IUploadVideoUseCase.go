// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	video "github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	mock "github.com/stretchr/testify/mock"
)

// IUploadVideoUseCase is an autogenerated mock type for the IUploadVideoUseCase type
type IUploadVideoUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: cmd
func (_m *IUploadVideoUseCase) Execute(cmd video.UploadVideoCommand) (*video.UploadVideoResponse, error) {
	ret := _m.Called(cmd)

	var r0 *video.UploadVideoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(video.UploadVideoCommand) (*video.UploadVideoResponse, error)); ok {
		return rf(cmd)
	}
	if rf, ok := ret.Get(0).(func(video.UploadVideoCommand) *video.UploadVideoResponse); ok {
		r0 = rf(cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*video.UploadVideoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(video.UploadVideoCommand) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIUploadVideoUseCase creates a new instance of IUploadVideoUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIUploadVideoUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IUploadVideoUseCase {
	mock := &IUploadVideoUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
