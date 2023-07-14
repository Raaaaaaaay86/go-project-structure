// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	entity "github.com/raaaaaaaay86/go-project-structure/domain/entity"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	repository "github.com/raaaaaaaay86/go-project-structure/domain/repository"
)

// VideoPostRepository is an autogenerated mock type for the VideoPostRepository type
type VideoPostRepository struct {
	mock.Mock
}

// CommitTx provides a mock function with given fields: tx
func (_m *VideoPostRepository) CommitTx(tx *gorm.DB) *gorm.DB {
	ret := _m.Called(tx)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(*gorm.DB) *gorm.DB); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Create provides a mock function with given fields: post
func (_m *VideoPostRepository) Create(post *entity.VideoPost) error {
	ret := _m.Called(post)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.VideoPost) error); ok {
		r0 = rf(post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindById provides a mock function with given fields: id
func (_m *VideoPostRepository) FindById(id uint) (*entity.VideoPost, error) {
	ret := _m.Called(id)

	var r0 *entity.VideoPost
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) (*entity.VideoPost, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uint) *entity.VideoPost); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.VideoPost)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForUpdate provides a mock function with given fields:
func (_m *VideoPostRepository) ForUpdate() repository.VideoPostRepository {
	ret := _m.Called()

	var r0 repository.VideoPostRepository
	if rf, ok := ret.Get(0).(func() repository.VideoPostRepository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.VideoPostRepository)
		}
	}

	return r0
}

// StartTx provides a mock function with given fields:
func (_m *VideoPostRepository) StartTx() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Update provides a mock function with given fields: post
func (_m *VideoPostRepository) Update(post *entity.VideoPost) error {
	ret := _m.Called(post)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.VideoPost) error); ok {
		r0 = rf(post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: tx
func (_m *VideoPostRepository) WithTx(tx *gorm.DB) repository.VideoPostRepository {
	ret := _m.Called(tx)

	var r0 repository.VideoPostRepository
	if rf, ok := ret.Get(0).(func(*gorm.DB) repository.VideoPostRepository); ok {
		r0 = rf(tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repository.VideoPostRepository)
		}
	}

	return r0
}

// NewVideoPostRepository creates a new instance of VideoPostRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVideoPostRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *VideoPostRepository {
	mock := &VideoPostRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
