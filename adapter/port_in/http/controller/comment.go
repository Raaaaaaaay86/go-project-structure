package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/domain/exception"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

type ICommentController interface {
	Create(ctx *gin.Context)
	Find(ctx *gin.Context)
	UserDelete(ctx *gin.Context)
	ForceDelete(ctx *gin.Context)
}

type CommentController struct {
	CreateUseCase              comment.IVideoCommentCreateUseCase
	FindByVideoIdUseCase       comment.IFindByVideoCQRS
	DeleteCommentUseCase       comment.IDeleteCommentUseCase
	ForceDeleteCommentUserCase comment.IForceDeleteCommentUseCase
	TracerProvider             *trace.TracerProvider
}

func NewCommentController(tracerProvider *trace.TracerProvider, createCommentUseCase comment.IVideoCommentCreateUseCase, findByVideoUseCase comment.IFindByVideoCQRS, deleteCommentUseCase comment.IDeleteCommentUseCase, forceDeleteUseCase comment.IForceDeleteCommentUseCase) *CommentController {
	return &CommentController{
		CreateUseCase:              createCommentUseCase,
		FindByVideoIdUseCase:       findByVideoUseCase,
		DeleteCommentUseCase:       deleteCommentUseCase,
		ForceDeleteCommentUserCase: forceDeleteUseCase,
		TracerProvider:             tracerProvider,
	}
}

func (c CommentController) Create(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(c.TracerProvider, ctx, pkg)
	defer span.End()

	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, res.Fail(exception.ErrUnauthorized.Error(), nil))
		tracing.RecordHttpError(span, http.StatusUnauthorized, exception.ErrUnauthorized)
		return
	}

	cmd := comment.CreateCommentCommand{
		AuthorId: token.(*jwt.CustomClaim).Uid,
	}
	err := ctx.ShouldBindJSON(&cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		tracing.RecordHttpError(span, http.StatusUnauthorized, err)
		return
	}

	response, err := c.CreateUseCase.Execute(newCtx, cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to create comment."))
		tracing.RecordHttpError(span, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}

func (c CommentController) Find(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(c.TracerProvider, ctx, pkg)
	defer span.End()

	var query comment.FindByVideoQuery
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		tracing.RecordHttpError(span, http.StatusBadRequest, err)
		return
	}

	response, err := c.FindByVideoIdUseCase.Execute(newCtx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to find comments."))
		tracing.RecordHttpError(span, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, res.Success(response))
}

func (c CommentController) UserDelete(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(c.TracerProvider, ctx, pkg)
	defer span.End()

	var command comment.DeleteCommentCommand
	err := ctx.ShouldBindJSON(&command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		tracing.RecordHttpError(span, http.StatusBadRequest, err)
		return
	}

	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, res.Fail("token is not found.", nil))
		tracing.RecordHttpError(span, http.StatusUnauthorized, err)
		return
	}

	claims := token.(*jwt.CustomClaim)
	command.ExecutorId = claims.Uid
	command.RoleIds = claims.Roles

	response, err := c.DeleteCommentUseCase.Execute(newCtx, command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), nil))
		tracing.RecordHttpError(span, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, res.Success(response))
}

func (c CommentController) ForceDelete(ctx *gin.Context) {
	newCtx, span := tracing.HttpSpanFactory(c.TracerProvider, ctx, pkg)
	defer span.End()

	var command comment.ForceDeleteCommentCommand
	err := ctx.ShouldBindJSON(&command)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		tracing.RecordHttpError(span, http.StatusBadRequest, err)
		return
	}

	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, res.Fail("token is not found.", nil))
		tracing.RecordHttpError(span, http.StatusUnauthorized, err)
		return
	}

	claims := token.(*jwt.CustomClaim)
	command.ExecutorId = claims.Uid
	command.RoleIds = claims.Roles

	response, err := c.ForceDeleteCommentUserCase.Execute(newCtx, command)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), nil))
		tracing.RecordHttpError(span, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, res.Success(response))
}
