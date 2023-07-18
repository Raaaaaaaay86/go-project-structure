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
	"github.com/raaaaaaaay86/go-project-structure/pkg/logger"
	"github.com/raaaaaaaay86/go-project-structure/pkg/mongo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
)

func main() {
	// Config
	config, err := configs.ReadYaml("./config/app_config.yaml")
	if err != nil {
		panic(err)
	}

	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	// TracerProviders
	appTracer := tracing.NewJaegerTracerProvider("application")
	ginTracer := tracing.NewJaegerTracerProvider("http")
	postgresTracer := tracing.NewJaegerTracerProvider("postgres")
	mongodbTracer := tracing.NewJaegerTracerProvider("mongodb")

	// Postgres Repositories
	db, err := gorm.NewPostgresConnection(config.Postgres)
	if err != nil {
		panic(err)
	}
	userRepository := postgres.NewUserRepository(postgresTracer, db)
	videoPostRepository := postgres.NewVideoPostRepository(postgresTracer, db)

	// MongoDB Repositories
	client, err := mongo.NewMongoDbConnection(config.MongoDB)
	if err != nil {
		panic(err)
	}
	videoCommentRepository := mongodb.NewVideoCommentRepository(mongodbTracer, client)

	// Helper Package
	fileUploader := bucket.NewLocalUploader(config.BucketPath.Raw)
	ffmpeg := convert.NewFfmpeg(config.BucketPath.Converted)

	// Use Case

	registerUseCase := auth.NewRegisterUserUseCase(appTracer, userRepository)
	loginUseCase := auth.NewLoginUseCase(appTracer, userRepository)
	uploadVideoUseCase := createVideo.NewUploadVideoUseCase(appTracer, fileUploader, ffmpeg)
	createVideoUseCase := createVideo.NewCreateVideoUseCase(appTracer, videoPostRepository, userRepository)
	createCommentUseCase := comment.NewCreateCommentUseCase(appTracer, videoCommentRepository, userRepository, videoPostRepository)
	findCommentByVideoUseCase := comment.NewFindByVideoUseCase(appTracer, videoCommentRepository)
	deleteCommentUseCase := comment.NewDeleteCommentUseCase(appTracer, videoCommentRepository)
	forceDeleteCommentUseCase := comment.NewForceDeleteCommentUseCase(appTracer, videoCommentRepository)

	// HTTP Server

	authController := controller.NewAuthenticationController(ginTracer, registerUseCase, loginUseCase)
	videoController := controller.NewVideoController(ginTracer, uploadVideoUseCase, createVideoUseCase)
	commentController := controller.NewCommentController(ginTracer, createCommentUseCase, findCommentByVideoUseCase, deleteCommentUseCase, forceDeleteCommentUseCase)

	httpPort := fmt.Sprintf(":%d", config.Http.Port)
	err = http.
		NewHttpServer(authController, videoController, commentController).
		Run(httpPort)
	if err != nil {
		panic(err)
	}
}
