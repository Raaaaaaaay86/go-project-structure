package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/pkg/jwt"
	"net/http"
)

func Token(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	if len(authorization) <= 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(authorization[7:])
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("token", token)

	ctx.Next()
}
