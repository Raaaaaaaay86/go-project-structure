package auth_test

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestLoginCQRS_Execute(t *testing.T) {
	type LoginTestCase struct {
		Username          string
		DecryptedPassword vo.DecryptedPassword
		UserRole          role.RoleId
		ExpectedErr       error
	}

	testCases := []LoginTestCase{
		{
			Username:          "user01",
			DecryptedPassword: "",
			UserRole:          role.User,
			ExpectedErr:       exception.ErrEmptyInput,
		},
		{
			Username:          "",
			DecryptedPassword: "anypassword",
			UserRole:          role.User,
			ExpectedErr:       exception.ErrEmptyInput,
		},
		{
			Username:          "user01",
			DecryptedPassword: "user01secret",
			UserRole:          role.User,
			ExpectedErr:       nil,
		},
		{
			Username:          "user01",
			DecryptedPassword: "wrongPassword",
			UserRole:          role.User,
			ExpectedErr:       exception.ErrWrongPassword,
		},
	}

	for i, testcase := range testCases {
		t.Logf("Start Test case[%d]", i)

		cmd := auth.LoginUserCommand{
			Username:          testcase.Username,
			DecryptedPassword: testcase.DecryptedPassword,
		}

		userRepository := mocks.NewUserRepository(t)
		userRepository.On("WithPreload").Return(userRepository).Maybe()
		expectedUser := entity.NewUser(cmd.Username, cmd.DecryptedPassword.Encrypt(), mock.Anything, *(entity.NewRole(testcase.UserRole)))
		switch testcase.ExpectedErr {
		case nil, exception.ErrWrongPassword:
			userRepository.On("FindByUsername", cmd.Username).Return(expectedUser, nil).Once()
		case exception.ErrUserNotFound:
			userRepository.On("FindByUsername", cmd.Username).Return(nil, gorm.ErrRecordNotFound).Once()
		case exception.ErrEmptyInput:
		}

		tracer := tracing.NewEmptyTracerProvider("application")
		response, err := auth.NewLoginUseCase(tracer, userRepository).Execute(context.TODO(), cmd)
		if err != nil {
			assert.ErrorIs(t, testcase.ExpectedErr, err)
			continue
		}

		assert.NotEmptyf(t, response.Token, "response.Token should not be empty")

		claim, err := jwt.Parse(response.Token)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Id, claim.Uid)

		expectedRoleIds := make([]role.RoleId, 0)
		for _, r := range expectedUser.Roles {
			expectedRoleIds = append(expectedRoleIds, r.Id)
		}
		for _, r := range claim.Roles {
			assert.Contains(t, expectedRoleIds, r)
		}
	}
}
