package route

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/internal/exception"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

var _ ICommentController = (*CommentController)(nil)

type ICommentController interface {
	Create(ctx *gin.Context)
	Find(ctx *gin.Context)
	UserDelete(ctx *gin.Context)
	ForceDelete(ctx *gin.Context)
}

type CommentController struct {
	CreateUseCase              comment.ICreateCommentUseCase
	FindByVideoIdUseCase       comment.IFindByVideoUseCase
	DeleteCommentUseCase       comment.IDeleteCommentUseCase
	ForceDeleteCommentUserCase comment.IForceDeleteCommentUseCase
	TracerProvider             *trace.TracerProvider
}

func NewCommentController(tracerProvider *trace.TracerProvider, createCommentUseCase comment.ICreateCommentUseCase, findByVideoUseCase comment.IFindByVideoUseCase, deleteCommentUseCase comment.IDeleteCommentUseCase, forceDeleteUseCase comment.IForceDeleteCommentUseCase) *CommentController {
	return &CommentController{
		CreateUseCase:              createCommentUseCase,
		FindByVideoIdUseCase:       findByVideoUseCase,
		DeleteCommentUseCase:       deleteCommentUseCase,
		ForceDeleteCommentUserCase: forceDeleteUseCase,
		TracerProvider:             tracerProvider,
	}
}

// @Summary	Create comment.
// @Tags		comment
// @Accept		json
// @Produce	json
// @Param		request	body		comment.CreateCommentCommand	true	"request body"
// @Success	200		{object}	comment.CreateCommentResponse
// @Router		/video/api/v1/comment/create [post]
// @Security	BearerAuth
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

// @Summary	Find comments by video id.
// @Tags		comment
// @Accept		json
// @Produce	json
// @Param		request	query		comment.FindByVideoQuery	true	"request body"
// @Success	200		{object}	comment.FindByVideoResponse
// @Router		/video/api/v1/comment/find [get]
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

// @Summary	Delete comment by normal user
// @Tags		comment
// @Accept		json
// @Produce	json
// @Param		request	body		comment.DeleteCommentCommand	true	"request body"
// @Success	200		{object}	comment.DeleteCommentResponse
// @Router		/video/api/v1/comment/delete [delete]
// @Security	BearerAuth
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

// @Summary	Force delete comment by admin role
// @Tags		comment
// @Accept		json
// @Produce	json
// @Param		request	body		comment.ForceDeleteCommentCommand	true	"request body"
// @Success	200		{object}	comment.ForceDeleteCommentResponse
// @Router		/video/api/v1/comment/admin/delete [delete]
// @Security	BearerAuth
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
