package auth

import (
	"context"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/internal/repository"
	"github.com/raaaaaaaay86/go-project-structure/internal/vo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"github.com/raaaaaaaay86/go-project-structure/pkg/validate"
)

var _ validate.Validator = (*LoginUserCommand)(nil)

type LoginUserCommand struct {
	Username          string               `json:"username,omitempty" example:"username01"`
	DecryptedPassword vo.DecryptedPassword `json:"password,omitempty" example:"password"`
}

func (c LoginUserCommand) Validate() error {
	validations := []struct {
		ValidatedResult bool
		Err             func() error
	}{
		{
			ValidatedResult: len(c.Username) > 0,
			Err:             func() error { return exception.NewInvalidInputError("username").ShouldNotEmpty() },
		},
		{
			ValidatedResult: len(c.DecryptedPassword) > 0,
			Err:             func() error { return exception.NewInvalidInputError("password").ShouldNotEmpty() },
		},
	}

	for _, validation := range validations {
		if !validation.ValidatedResult {
			return validation.Err()
		}
	}

	return nil
}

type LoginUserResponse struct {
	Token string `json:"token,omitempty"`
}

type ILoginUserUseCase interface {
	Execute(ctx context.Context, cmd LoginUserCommand) (*LoginUserResponse, error)
}

var _ ILoginUserUseCase = (*LoginUserUseCase)(nil)

type LoginUserUseCase struct {
	UserRepository repository.UserRepository
	TracerProvider tracing.ApplicationTracer
}

func NewLoginUseCase(tracer tracing.ApplicationTracer, userRepository repository.UserRepository) *LoginUserUseCase {
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
