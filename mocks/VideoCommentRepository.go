// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	entity "github.com/raaaaaaaay86/go-project-structure/domain/entity"
	mock "github.com/stretchr/testify/mock"
)

// VideoCommentRepository is an autogenerated mock type for the VideoCommentRepository type
type VideoCommentRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: comment
func (_m *VideoCommentRepository) Create(comment *entity.VideoComment) error {
	ret := _m.Called(comment)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.VideoComment) error); ok {
		r0 = rf(comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByVideoId provides a mock function with given fields: videoId
func (_m *VideoCommentRepository) FindByVideoId(videoId uint) ([]*entity.VideoComment, error) {
	ret := _m.Called(videoId)

	var r0 []*entity.VideoComment
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) ([]*entity.VideoComment, error)); ok {
		return rf(videoId)
	}
	if rf, ok := ret.Get(0).(func(uint) []*entity.VideoComment); ok {
		r0 = rf(videoId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.VideoComment)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(videoId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewVideoCommentRepository creates a new instance of VideoCommentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVideoCommentRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *VideoCommentRepository {
	mock := &VideoCommentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
