package auth

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/entity"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo/enum/role"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ IRegisterUserUseCase = (*RegisterUserUseCase)(nil)

type RegisterUserCommand struct {
	Username          string               `json:"username,omitempty" example:"username01"`
	DecryptedPassword vo.DecryptedPassword `json:"password,omitempty" example:"password"`
	Email             vo.Email             `json:"email,omitempty" example:"example@gmail.com"`
}

func (c RegisterUserCommand) Validate() error {
	if c.Username == "" || c.DecryptedPassword == "" || c.Email == "" {
		return exception.ErrEmptyInput
	}
	if err := c.DecryptedPassword.Validate(); err != nil {
		return err
	}
	if err := c.Email.Validate(); err != nil {
		return err
	}
	return nil
}

type RegisterUserResponse struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type IRegisterUserUseCase interface {
	Execute(ctx context.Context, command RegisterUserCommand) (*RegisterUserResponse, error)
}

type RegisterUserUseCase struct {
	userRepository repository.UserRepository
	TracerProvider *trace.TracerProvider
}

func NewRegisterUserUseCase(tracerProvider *trace.TracerProvider, userRepository repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository: userRepository,
		TracerProvider: tracerProvider,
	}
}

func (u RegisterUserUseCase) Execute(ctx context.Context, cmd RegisterUserCommand) (*RegisterUserResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(u.TracerProvider, ctx, pkg, "RegisterUserUseCase.Execute")
	defer span.End()

	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(cmd.Username, cmd.DecryptedPassword.Encrypt(), cmd.Email, *(entity.NewRole(role.User)))

	err = u.userRepository.Create(newCtx, user)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &RegisterUserResponse{ID: user.Id, Username: user.Username}, nil
}
