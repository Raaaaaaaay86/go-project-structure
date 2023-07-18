package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

type IVideoController interface {
	Upload(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type VideoController struct {
	UploadVideoUseCase video.IUploadVideoUseCase
	CreateVideoUseCase video.IVideoCreateUseCase
	TracerProvider     *trace.TracerProvider
}

func NewVideoController(tracerProvider *trace.TracerProvider, upload video.IUploadVideoUseCase, create video.IVideoCreateUseCase) *VideoController {
	return &VideoController{
		UploadVideoUseCase: upload,
		CreateVideoUseCase: create,
		TracerProvider:     tracerProvider,
	}
}

// @Summary	Upload video.
// @Tags		video
// @Accept		json
// @Produce	json
// @Param		request	body		video.UploadVideoCommand	true	"request body"
// @Success	200		{object}	video.UploadVideoResponse
// @Router		/video/api/v1/video/upload [post]
// @Security	BearerAuth
func (v VideoController) Upload(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(v.TracerProvider, ctx, pkg)
	defer span.End()

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
	response, err := v.UploadVideoUseCase.Execute(newCtx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}

// @Summary	Create video post. This API is called after the video upload. Should bring the uploaded video uuid in the request.
// @Tags		video
// @Accept		json
// @Produce	json
// @Param		request	body		video.CreateVideoCommand	true	"request body"
// @Success	200		{object}	video.CreateVideoResponse
// @Router		/video/api/v1/create [post]
// @Security	BearerAuth
func (v VideoController) Create(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(v.TracerProvider, ctx, pkg)
	defer span.End()

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

	response, err := v.CreateVideoUseCase.Execute(newCtx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to create video post."))
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}
