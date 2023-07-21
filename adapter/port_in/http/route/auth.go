package route

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

var _ IAuthenticationController = (*AuthenticationController)(nil)

type IAuthenticationController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthenticationController struct {
	RegisterUseCase auth.IRegisterUserUseCase
	LoginUseCase    auth.ILoginUserResponse
	TracerProvider  *trace.TracerProvider
}

func NewAuthenticationController(tracerProvider *trace.TracerProvider, registerUseCase auth.IRegisterUserUseCase, loginUseCase auth.ILoginUserResponse) AuthenticationController {
	return AuthenticationController{
		RegisterUseCase: registerUseCase,
		LoginUseCase:    loginUseCase,
		TracerProvider:  tracerProvider,
	}
}

// @Summary	Register new user.
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		request	body		auth.RegisterUserCommand	true	"request format"
// @Success	200		{object}	auth.RegisterUserResponse
// @Router		/v1/auth/register [post]
// @Security	BearerAuth
func (a AuthenticationController) Register(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(a.TracerProvider, ctx, pkg)
	defer span.End()

	var command auth.RegisterUserCommand
	err := ctx.ShouldBindJSON(&command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), nil))
		tracing.RecordHttpError(span, http.StatusBadRequest, err)
		return
	}

	response, err := a.RegisterUseCase.Execute(newCtx, command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), nil))
		tracing.RecordHttpError(span, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}

// @Summary	Login and get token.
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		request	body		auth.RegisterUserCommand	true	"request format"
// @Success	200		{object}	auth.RegisterUserResponse
// @Router		/v1/auth/login [post]
// @Security	BearerAuth
func (a AuthenticationController) Login(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(a.TracerProvider, ctx, pkg)
	defer span.End()

	var command auth.LoginUserCommand
	err := ctx.ShouldBindJSON(&command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), nil))
		tracing.RecordHttpError(span, http.StatusBadRequest, err)
		return
	}

	response, err := a.LoginUseCase.Execute(newCtx, command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), nil))
		tracing.RecordHttpError(span, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, res.Success(response))
}
