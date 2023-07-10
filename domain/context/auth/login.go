package auth

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
)

type LoginUserCommand struct {
	Username          string               `json:"username,omitempty" example:"username01"`
	DecryptedPassword vo.DecryptedPassword `json:"password,omitempty" example:"password"`
}

func (c LoginUserCommand) Validate() error {
	if len(c.Username) <= 0 || len(c.DecryptedPassword) <= 0 {
		return exception.ErrEmptyInput
	}
	return nil
}

type LoginUserResponse struct {
	Token string `json:"token,omitempty"`
}

type ILoginUserResponse interface {
	Execute(cmd LoginUserCommand) (*LoginUserResponse, error)
}

type LoginUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewLoginUseCase(userRepository repository.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		UserRepository: userRepository,
	}
}

func (c LoginUserUseCase) Execute(cmd LoginUserCommand) (*LoginUserResponse, error) {
	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	user, err := c.UserRepository.WithPreload().FindByUsername(cmd.Username)
	if err != nil {
		return nil, exception.ErrUserNotFound
	}

	if !user.Password.Compare(cmd.DecryptedPassword) {
		return nil, exception.ErrWrongPassword
	}

	tokenString, err := jwt.Generate(*user)
	if err != nil {
		return nil, exception.ErrTokenGenerationFailed
	}

	return &LoginUserResponse{Token: tokenString}, nil
}
