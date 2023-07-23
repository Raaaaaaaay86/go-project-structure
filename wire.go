//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http/route"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_out/repository"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg"
)

type Application struct {
	HttpServer *gin.Engine
}

func App() (*Application, error) {
	wire.Build(
		route.RouteSet,
		http.HttpServerSet,
		repository.RepositorySet,
		auth.ContextAuthSet,
		comment.ContextCommentSet,
		video.ContextVideoSet,
		pkg.PkgSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
