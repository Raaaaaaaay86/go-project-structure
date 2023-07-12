package main

import (
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http/controller"
	mongodb "github.com/raaaaaaaay86/go-project-structure/adapter/port_out/mongo"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_out/postgres"
	_ "github.com/raaaaaaaay86/go-project-structure/docs" //nolint:typecheck
	"github.com/raaaaaaaay86/go-project-structure/domain/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	createVideo "github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	convert "github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/gorm"
	"github.com/raaaaaaaay86/go-project-structure/pkg/mongo"
)

func main() {
	// Config
	config, err := configs.ReadYaml("./config/app_config.yaml")
	if err != nil {
		panic(err)
	}

	// Postgres Repositories
	db, err := gorm.NewPostgresConnection(config.Postgres)
	if err != nil {
		panic(err)
	}
	userRepository := postgres.NewUserRepository(db)
	videoPostRepository := postgres.NewVideoPostRepository(db)

	// MongoDB Repositories
	client, err := mongo.NewMongoDbConnection(config.MongoDB)
	if err != nil {
		panic(err)
	}
	commentsCollection := client.Database("video").Collection("comments")
	videoCommentRepository := mongodb.NewVideoCommentRepository(commentsCollection)

	// Helper Package
	fileUploader := bucket.NewLocalUploader(config.BucketPath.Raw)
	ffmpeg := convert.NewFfmpeg(config.BucketPath.Converted)

	// Use Case
	registerUseCase := auth.NewRegisterUserUseCase(userRepository)
	loginUseCase := auth.NewLoginUseCase(userRepository)
	uploadVideoUseCase := createVideo.NewUploadVideoUseCase(fileUploader, ffmpeg)
	createVideoUseCase := createVideo.NewCreateVideoUseCase(videoPostRepository, userRepository)
	createCommentUseCase := comment.NewCreateCommentUseCase(videoCommentRepository, userRepository, videoPostRepository)
	findCommentByVideoUseCase := comment.NewFindByVideoUseCase(videoCommentRepository)

	// HTTP Server
	authController := controller.NewAuthenticationController(registerUseCase, loginUseCase)
	videoController := controller.NewVideoController(uploadVideoUseCase, createVideoUseCase)
	commentController := controller.NewCommentController(createCommentUseCase, findCommentByVideoUseCase)

	httpPort := fmt.Sprintf(":%d", config.Http.Port)
	err = http.
		NewHttpServer(authController, videoController, commentController).
		Run(httpPort)
	if err != nil {
		panic(err)
	}
}
