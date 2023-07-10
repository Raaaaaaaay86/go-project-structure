package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/res"
	"net/http"
)

type ICommentController interface {
	Create(ctx *gin.Context)
	Find(ctx *gin.Context)
}

type CommentController struct {
	CreateUseCase        comment.IVideoCommentCreateUseCase
	FindByVideoIdUseCase comment.IFindByVideoCQRS
}

func NewCommentController(createCommentUseCase comment.IVideoCommentCreateUseCase, findByVideoUseCase comment.IFindByVideoCQRS) *CommentController {
	return &CommentController{
		CreateUseCase:        createCommentUseCase,
		FindByVideoIdUseCase: findByVideoUseCase,
	}
}

func (c CommentController) Create(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, res.Fail("token is not found.", nil))
		return
	}

	cmd := comment.CreateCommentCommand{
		AuthorId: token.(*jwt.CustomClaim).Uid,
	}
	err := ctx.ShouldBindJSON(&cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		return
	}

	response, err := c.CreateUseCase.Execute(cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to create comment."))
		return
	}

	ctx.JSON(http.StatusCreated, res.Success(response))
}

func (c CommentController) Find(ctx *gin.Context) {
	var query comment.FindByVideoQuery
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, res.Fail(err.Error(), "invalid request body."))
		return
	}

	response, err := c.FindByVideoIdUseCase.Execute(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Fail(err.Error(), "unable to find comments."))
		return
	}

	ctx.JSON(http.StatusOK, res.Success(response))
}
