package main

import (
	"context"
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/http/route"
	mongodb "github.com/raaaaaaaay86/go-project-structure/adapter/port_out/repository"
	_ "github.com/raaaaaaaay86/go-project-structure/docs" //nolint:typecheck
	"github.com/raaaaaaaay86/go-project-structure/domain/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/domain/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	convert "github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/gorm"
	"github.com/raaaaaaaay86/go-project-structure/pkg/graphdb"
	"github.com/raaaaaaaay86/go-project-structure/pkg/logger"
	"github.com/raaaaaaaay86/go-project-structure/pkg/mongo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
)

func main() {
	ctx := context.Background()
	// Config
	config, err := configs.ReadYaml("./config/app_config.yaml")
	if err != nil {
		panic(err)
	}

	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync() //nolint:errcheck

	// Tracer
	exporter, err := tracing.NewJaegerExporter(config.Jaeger.Endpoint)
	if err != nil {
		zapLogger.Error("jaeger exporter error", err)
	}

	appTracer, err := tracing.NewJaegerTracerProvider("application", exporter)
	if err != nil {
		zapLogger.Error("init application tracer failed.", err)
	}

	httpTracer, err := tracing.NewJaegerTracerProvider("http", exporter)
	if err != nil {
		zapLogger.Error("init http tracer failed", err)
	}

	repositoryTracer, err := tracing.NewJaegerTracerProvider("repository", exporter)
	if err != nil {
		zapLogger.Error("init repository tracer failed", err)
	}

	// Neo4j Driver
	neo4jDriver, err := graphdb.NewNeo4j(ctx, config.Neo4j)
	if err != nil {
		panic(err)
	}
	defer neo4jDriver.Close(ctx)

	// Postgres Repositories
	db, err := gorm.NewPostgresConnection(config.Postgres)
	if err != nil {
		panic(err)
	}
	userRepository := mongodb.NewUserRepository(repositoryTracer, db)
	videoPostRepository := mongodb.NewVideoPostRepository(repositoryTracer, db, neo4jDriver)

	// MongoDB Repositories
	client, err := mongo.NewMongoDbConnection(config.MongoDB)
	if err != nil {
		panic(err)
	}
	videoCommentRepository := mongodb.NewVideoCommentRepository(repositoryTracer, client)

	// Helper Package
	fileUploader := bucket.NewLocalUploader(config.BucketPath.Raw)
	ffmpeg := convert.NewFfmpeg(config.BucketPath.Converted)

	// Use Case

	registerUseCase := auth.NewRegisterUserUseCase(appTracer, userRepository)
	loginUseCase := auth.NewLoginUseCase(appTracer, userRepository)
	uploadVideoUseCase := video.NewUploadVideoUseCase(appTracer, fileUploader, ffmpeg)
	createVideoUseCase := video.NewCreateVideoUseCase(appTracer, videoPostRepository, userRepository)
	createCommentUseCase := comment.NewCreateCommentUseCase(appTracer, videoCommentRepository, userRepository, videoPostRepository)
	findCommentByVideoUseCase := comment.NewFindByVideoUseCase(appTracer, videoCommentRepository)
	deleteCommentUseCase := comment.NewDeleteCommentUseCase(appTracer, videoCommentRepository)
	forceDeleteCommentUseCase := comment.NewForceDeleteCommentUseCase(appTracer, videoCommentRepository)
	likeVideoUseCase := video.NewLikeVideoUseCase(appTracer, videoPostRepository)
	unlikeVideoUseCase := video.NewUnLikeVideoUseCase(appTracer, videoPostRepository)

	// HTTP Server

	authController := route.NewAuthenticationController(httpTracer, registerUseCase, loginUseCase)
	videoController := route.NewVideoController(httpTracer, uploadVideoUseCase, createVideoUseCase, likeVideoUseCase, unlikeVideoUseCase)
	commentController := route.NewCommentController(httpTracer, createCommentUseCase, findCommentByVideoUseCase, deleteCommentUseCase, forceDeleteCommentUseCase)

	httpPort := fmt.Sprintf(":%d", config.Http.Port)
	err = http.
		NewHttpServer(authController, videoController, commentController).
		Run(httpPort)
	if err != nil {
		panic(err)
	}
}
