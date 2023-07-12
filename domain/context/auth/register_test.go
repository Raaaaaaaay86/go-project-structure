package auth_test

import (
	r "github.com/raaaaaaaay86/go-project-structure/domain/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	type RegisterTestCase struct {
		Username          string
		RoleId            role.RoleId
		DecryptedPassword vo.DecryptedPassword
		Email             vo.Email
		ExceptedErr       error
	}

	userRepo := mocks.NewUserRepository(t)
	register := r.NewRegisterUserUseCase(userRepo)

	testCases := []RegisterTestCase{
		{
			Username:          "user01",
			DecryptedPassword: "user01secret",
			Email:             "user01@email.com",
			ExceptedErr:       nil,
		},
		{
			Username:          "user01",
			DecryptedPassword: "user01secret",
			Email:             "user01email.com",
			ExceptedErr:       exception.ErrInvalidEmail,
		},
		{
			Username:          "user01",
			DecryptedPassword: "",
			Email:             "user01@email.com",
			ExceptedErr:       exception.ErrEmptyInput,
		},
	}

	for i, tc := range testCases {
		t.Logf("Start Test case[%d]", i)

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
			userRepo.On("Create", mock.IsType(newUser)).Return(nil).Once()
		}

		response, err := register.Execute(cmd)
		if err != nil {
			assert.ErrorIs(t, tc.ExceptedErr, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, newUser.Username, response.Username)
	}
}
