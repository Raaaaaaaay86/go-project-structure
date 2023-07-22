package auth_test

import (
	"context"
	r "github.com/raaaaaaaay86/go-project-structure/internal/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	type RegisterTestCase struct {
		TestDescription   string
		Username          string
		RoleId            role.RoleId
		DecryptedPassword vo.DecryptedPassword
		Email             vo.Email
		ExceptedErr       error
	}

	userRepo := mocks.NewUserRepository(t)
	register := r.NewRegisterUserUseCase(tracing.NewEmptyTracerProvider(), userRepo)

	testCases := []RegisterTestCase{
		{
			TestDescription:   "Success Register",
			Username:          "user01",
			DecryptedPassword: "correctPassword",
			Email:             "user01@email.com",
			ExceptedErr:       nil,
		},
		{
			TestDescription:   "Failed by empty invalid email",
			Username:          "user01",
			DecryptedPassword: "correctPassword",
			Email:             "user01email.com",
			ExceptedErr:       exception.ErrInvalidEmail,
		},
		{
			TestDescription:   "Failed by empty password",
			Username:          "user01",
			DecryptedPassword: "",
			Email:             "user01@email.com",
			ExceptedErr:       exception.NewInvalidInputError("password").ShouldNotEmpty(),
		},
		{
			TestDescription:   "Failed by empty username",
			Username:          "",
			DecryptedPassword: "correctPassword",
			Email:             "user01@email.com",
			ExceptedErr:       exception.NewInvalidInputError("username").ShouldNotEmpty(),
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d] - %s", i, tc.TestDescription)

		cmd := r.RegisterUserCommand{
			Username:          tc.Username,
			DecryptedPassword: tc.DecryptedPassword,
			Email:             tc.Email,
		}

		newUser := entity.NewUser(
			cmd.Username,
			cmd.DecryptedPassword.Encrypt(),
			cmd.Email,
			*(entity.NewRole(role.User)),
		)

		if tc.ExceptedErr == nil {
			userRepo.On("Create", mock.Anything, mock.IsType(newUser)).Return(nil).Once()
		}

		response, err := register.Execute(context.TODO(), cmd)
		if err != nil || tc.ExceptedErr != nil {
			assert.ErrorIs(t, tc.ExceptedErr, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, newUser.Username, response.Username)
	}
}
