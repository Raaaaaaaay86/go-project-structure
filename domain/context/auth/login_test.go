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
		TestDescription   string
		Username          string
		DecryptedPassword vo.DecryptedPassword
		UserRole          role.RoleId
		ExpectedErr       error
	}

	testCases := []LoginTestCase{
		{
			TestDescription:   "Failed by empty password",
			Username:          "user01",
			DecryptedPassword: "",
			UserRole:          role.User,
			ExpectedErr:       exception.ErrEmptyInput,
		},
		{
			TestDescription:   "Failed by empty username",
			Username:          "",
			DecryptedPassword: "correctPassword",
			UserRole:          role.User,
			ExpectedErr:       exception.ErrEmptyInput,
		},
		{
			TestDescription:   "Failed by wrong password",
			Username:          "user01",
			DecryptedPassword: "wrongPassword",
			UserRole:          role.User,
			ExpectedErr:       exception.ErrWrongPassword,
		},
		{
			TestDescription:   "Success Login",
			Username:          "user01",
			DecryptedPassword: "correctPassword",
			UserRole:          role.User,
			ExpectedErr:       nil,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d] - %s", i, tc.TestDescription)

		cmd := auth.LoginUserCommand{
			Username:          tc.Username,
			DecryptedPassword: tc.DecryptedPassword,
		}

		ctx := context.Background()
		userRepository := mocks.NewUserRepository(t)
		userRepository.On("WithPreload").Return(userRepository).Maybe()
		expectedUser := entity.NewUser(cmd.Username, vo.DecryptedPassword("correctPassword").Encrypt(), mock.Anything, *(entity.NewRole(tc.UserRole)))
		switch tc.ExpectedErr {
		case nil, exception.ErrWrongPassword:
			userRepository.On("FindByUsername", mock.Anything, cmd.Username).Return(expectedUser, nil).Once()
		case exception.ErrUserNotFound:
			userRepository.On("FindByUsername", mock.Anything, cmd.Username).Return(nil, gorm.ErrRecordNotFound).Once()
		case exception.ErrEmptyInput:
		}

		tracer := tracing.NewEmptyTracerProvider()
		response, err := auth.NewLoginUseCase(tracer, userRepository).Execute(ctx, cmd)
		if err != nil || tc.ExpectedErr != nil {
			assert.ErrorIs(t, tc.ExpectedErr, err)
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
