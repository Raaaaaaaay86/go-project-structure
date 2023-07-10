package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"net/http"
)

type IVideoController interface {
	Upload(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type VideoController struct {
	UploadVideoUseCase video.IUploadVideoUseCase
	CreateVideoUseCase video.IVideoCreateUseCase
}

func NewVideoController(upload video.IUploadVideoUseCase, create video.IVideoCreateUseCase) *VideoController {
	return &VideoController{
		UploadVideoUseCase: upload,
		CreateVideoUseCase: create,
	}
}

func (v VideoController) Upload(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("v")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "file is not included in request."))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to open file."))
		return
	}
	defer file.Close()

	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, res.Fail("token is not found.", nil))
		return
	}

	cmd := video.UploadVideoCommand{
		File:       file,
		FileName:   fileHeader.Filename,
		UploaderId: (token.(*jwt.CustomClaim)).Uid,
	}
	response, err := v.UploadVideoUseCase.Execute(cmd)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}

func (v VideoController) Create(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, res.Fail("token is not found.", nil))
		return
	}

	cmd := video.CreateVideoCommand{
		AuthorId: (token.(*jwt.CustomClaim)).Uid,
	}
	err := ctx.ShouldBindJSON(&cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		return
	}

	response, err := v.CreateVideoUseCase.Execute(cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to create video post."))
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}
