package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"net/http"
)

type IAuthenticationController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthenticationController struct {
	RegisterUseCase auth.IRegisterUserUseCase
	LoginUseCase    auth.ILoginUserResponse
}

func NewAuthenticationController(registerUseCase auth.IRegisterUserUseCase, loginUseCase auth.ILoginUserResponse) AuthenticationController {
	return AuthenticationController{
		RegisterUseCase: registerUseCase,
		LoginUseCase:    loginUseCase,
	}
}

// Register
//
//		@Summary					Register new user.
//		@Tags						auth
//		@Accept						json
//		@Produce					json
//		@Param						request	body		register.RegisterUserCommand	true	"request format"
//		@Success					200		{object}	register.RegisterUserResponse
//		@Router						/v1/auth/register [post]
//	 @Security BearerAuth
func (a AuthenticationController) Register(ctx *gin.Context) {
	var command auth.RegisterUserCommand
	err := ctx.ShouldBindJSON(&command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), nil))
		return
	}

	response, err := a.RegisterUseCase.Execute(command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}

// Login
//
//		@Summary					Login and get token.
//		@Tags						auth
//		@Accept						json
//		@Produce					json
//		@Param						request	body		login.RegisterUserCommand	true	"request format"
//		@Success					200		{object}	login.RegisterUserResponse
//		@Router						/v1/auth/login [post]
//	 @Security BearerAuth
func (a AuthenticationController) Login(ctx *gin.Context) {
	var command auth.LoginUserCommand
	err := ctx.ShouldBindJSON(&command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), nil))
		return
	}

	response, err := a.LoginUseCase.Execute(command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, res.Success(response))
}
