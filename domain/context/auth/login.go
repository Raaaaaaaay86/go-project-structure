package auth

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/domain/repository"
	"github.com/raaaaaaaay86/go-project-structure/domain/vo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
	"go.opentelemetry.io/otel/sdk/trace"
)

var _ ILoginUserResponse = (*LoginUserUseCase)(nil)

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
	Execute(ctx context.Context, cmd LoginUserCommand) (*LoginUserResponse, error)
}

type LoginUserUseCase struct {
	UserRepository repository.UserRepository
	TracerProvider *trace.TracerProvider
}

func NewLoginUseCase(tracer *trace.TracerProvider, userRepository repository.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		UserRepository: userRepository,
		TracerProvider: tracer,
	}
}

func (c LoginUserUseCase) Execute(ctx context.Context, cmd LoginUserCommand) (*LoginUserResponse, error) {
	newCtx, span := tracing.ApplicationSpanFactory(c.TracerProvider, ctx, pkg, "LoginUserUseCase.Execute")
	defer span.End()

	err := validate.Do(cmd)
	if err != nil {
		return nil, err
	}

	user, err := c.UserRepository.FindByUsername(newCtx, cmd.Username)
	if err != nil {
		span.RecordError(err)
		return nil, exception.ErrUserNotFound
	}

	if !user.Password.Compare(cmd.DecryptedPassword) {
		return nil, exception.ErrWrongPassword
	}

	tokenString, err := jwt.Generate(*user)
	if err != nil {
		return nil, exception.ErrTokenGenerationFailed
	}

	res := &LoginUserResponse{
		Token: tokenString,
	}

	return res, nil
}
