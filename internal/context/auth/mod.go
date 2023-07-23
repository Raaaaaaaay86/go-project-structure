package auth

import "github.com/google/wire"

const pkg = "internal.context.auth"

var ContextAuthSet = wire.NewSet(
	NewLoginUseCase,
	wire.Bind(new(ILoginUserUseCase), new(*LoginUserUseCase)),
	NewRegisterUserUseCase,
	wire.Bind(new(IRegisterUserUseCase), new(*RegisterUserUseCase)),
)
