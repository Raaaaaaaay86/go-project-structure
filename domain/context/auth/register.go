package auth

import (
	"github.com/raaaaaaaay86/go-project-structure/domain/entity"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
)

type RegisterUserCommand struct {
	Username          string               `json:"username,omitempty" example:"username01"`
	DecryptedPassword vo.DecryptedPassword `json:"password,omitempty" example:"password"`
	Email             vo.Email             `json:"email,omitempty" example:"example@gmail.com"`
}

type RegisterUserResponse struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type IRegisterUserUseCase interface {
	Execute(command RegisterUserCommand) (*RegisterUserResponse, error)
}

type RegisterUserUseCase struct {
	userRepository repository.UserRepository
}

func NewRegisterUserUseCase(userRepository repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepository: userRepository}
}

func (u RegisterUserUseCase) Execute(cmd RegisterUserCommand) (*RegisterUserResponse, error) {
	err := validate.Do(cmd.Email, cmd.DecryptedPassword)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(cmd.Username, cmd.DecryptedPassword.Encrypt(), cmd.Email, *(entity.NewRole(role.User, role.User.Code())))

	err = u.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return &RegisterUserResponse{ID: user.Id, Username: user.Username}, nil
}
