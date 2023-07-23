// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http/route"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_out/repository"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	"github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/gorm"
	"github.com/raaaaaaaay86/go-project-structure/pkg/graphdb"
	"github.com/raaaaaaaay86/go-project-structure/pkg/mongo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
)

import (
	_ "github.com/raaaaaaaay86/go-project-structure/docs"
)

// Injectors from wire.go:

func App() (*Application, error) {
	yamlConfig, err := configs.YamlConfigProvider()
	if err != nil {
		return nil, err
	}
	exporter, err := tracing.NewJaegerExporter(yamlConfig)
	if err != nil {
		return nil, err
	}
	httpTracer, err := tracing.NewHttpTracer(exporter)
	if err != nil {
		return nil, err
	}
	applicationTracer, err := tracing.NewApplicationTracer(exporter)
	if err != nil {
		return nil, err
	}
	repositoryTracer, err := tracing.NewRepositoryTracer(exporter)
	if err != nil {
		return nil, err
	}
	db, err := gorm.NewPostgresConnection(yamlConfig)
	if err != nil {
		return nil, err
	}
	userRepository := repository.NewUserRepository(repositoryTracer, db)
	registerUserUseCase := auth.NewRegisterUserUseCase(applicationTracer, userRepository)
	loginUserUseCase := auth.NewLoginUseCase(applicationTracer, userRepository)
	authenticationController := route.NewAuthenticationController(httpTracer, registerUserUseCase, loginUserUseCase)
	localUploader := bucket.NewLocalUploader(yamlConfig)
	ffmpeg := convert.NewFfmpeg(yamlConfig)
	uploadVideoUseCase := video.NewUploadVideoUseCase(applicationTracer, localUploader, ffmpeg)
	driverWithContext, err := graphdb.NewNeo4j(yamlConfig)
	if err != nil {
		return nil, err
	}
	videoPostRepository := repository.NewVideoPostRepository(repositoryTracer, db, driverWithContext)
	createVideoUseCase := video.NewCreateVideoUseCase(applicationTracer, videoPostRepository, userRepository)
	likeVideoUseCase := video.NewLikeVideoUseCase(applicationTracer, videoPostRepository)
	unLikeVideoUseCase := video.NewUnLikeVideoUseCase(applicationTracer, videoPostRepository)
	videoController := route.NewVideoController(httpTracer, uploadVideoUseCase, createVideoUseCase, likeVideoUseCase, unLikeVideoUseCase)
	client, err := mongo.NewMongoDbConnection(yamlConfig)
	if err != nil {
		return nil, err
	}
	videoCommentRepository := repository.NewVideoCommentRepository(repositoryTracer, client)
	createCommentUseCase := comment.NewCreateCommentUseCase(applicationTracer, videoCommentRepository, userRepository, videoPostRepository)
	findByVideoUseCase := comment.NewFindByVideoUseCase(applicationTracer, videoCommentRepository)
	deleteCommentUseCase := comment.NewDeleteCommentUseCase(applicationTracer, videoCommentRepository)
	forceDeleteCommentUseCase := comment.NewForceDeleteCommentUseCase(applicationTracer, videoCommentRepository)
	commentController := route.NewCommentController(httpTracer, createCommentUseCase, findByVideoUseCase, deleteCommentUseCase, forceDeleteCommentUseCase)
	engine := http.NewHttpServer(authenticationController, videoController, commentController)
	application := &Application{
		HttpServer: engine,
	}
	return application, nil
}

// wire.go:

type Application struct {
	HttpServer *gin.Engine
}
