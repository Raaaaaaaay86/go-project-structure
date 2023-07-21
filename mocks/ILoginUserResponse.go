// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	context "context"

	auth "github.com/raaaaaaaay86/go-project-structure/internal/context/auth"

	mock "github.com/stretchr/testify/mock"
)

// ILoginUserResponse is an autogenerated mock type for the ILoginUserResponse type
type ILoginUserResponse struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, cmd
func (_m *ILoginUserResponse) Execute(ctx context.Context, cmd auth.LoginUserCommand) (*auth.LoginUserResponse, error) {
	ret := _m.Called(ctx, cmd)

	var r0 *auth.LoginUserResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.LoginUserCommand) (*auth.LoginUserResponse, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auth.LoginUserCommand) *auth.LoginUserResponse); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auth.LoginUserResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auth.LoginUserCommand) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewILoginUserResponse creates a new instance of ILoginUserResponse. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewILoginUserResponse(t interface {
	mock.TestingT
	Cleanup(func())
}) *ILoginUserResponse {
	mock := &ILoginUserResponse{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
