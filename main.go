package main

import (
	"context"
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/httplayer"
	"github.com/raaaaaaaay86/go-project-structure/adapter/port_in/httplayer/route"
	mongodb "github.com/raaaaaaaay86/go-project-structure/adapter/port_out/repository"
	_ "github.com/raaaaaaaay86/go-project-structure/docs" //nolint:typecheck
	"github.com/raaaaaaaay86/go-project-structure/internal/context/auth"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/comment"
	"github.com/raaaaaaaay86/go-project-structure/internal/context/media/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	convert "github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/gorm"
	"github.com/raaaaaaaay86/go-project-structure/pkg/graphdb"
	"github.com/raaaaaaaay86/go-project-structure/pkg/logger"
	"github.com/raaaaaaaay86/go-project-structure/pkg/mongo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func Run(ctx context.Context) error {
	// Config
	config, err := configs.ReadYaml("./config/app_config.yaml")
	if err != nil {
		return err
	}

	// Logger
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		return err
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
		zapLogger.Error("init httplayer tracer failed", err)
	}

	repositoryTracer, err := tracing.NewJaegerTracerProvider("repository", exporter)
	if err != nil {
		zapLogger.Error("init repository tracer failed", err)
	}

	gormTracer, err := tracing.NewJaegerTracerProvider("gorm", exporter)
	if err != nil {
		zapLogger.Error("init gorm tracer failed", err)
	}

	// Neo4j Driver
	neo4jDriver, err := graphdb.NewNeo4j(ctx, config.Neo4j)
	if err != nil {
		return err
	}
	defer neo4jDriver.Close(ctx)

	// Postgres Repositories
	gormDB, err := gorm.NewPostgresConnection(config.Postgres, gormTracer)
	if err != nil {
		neo4jDriver.Close(ctx)
		return err
	}
	userRepository := mongodb.NewUserRepository(repositoryTracer, gormDB)
	videoPostRepository := mongodb.NewVideoPostRepository(repositoryTracer, gormDB, neo4jDriver)

	// MongoDB Repositories
	mongoClient, err := mongo.NewMongoDbConnection(config.MongoDB)
	if err != nil {
		neo4jDriver.Close(ctx)
		mongoClient.Disconnect(ctx) //nolint:errcheck
		return err
	}
	videoCommentRepository := mongodb.NewVideoCommentRepository(repositoryTracer, mongoClient)

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

	// Start HTTP Server
	port := fmt.Sprintf(":%d", config.Http.Port)
	server := &http.Server{
		Addr:    port,
		Handler: httplayer.NewHttpServer(authController, videoController, commentController),
	}
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// HTTP Server Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	exitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	zapLogger.Warnw("server shutdown in 5 seconds", "event", "SERVER_SHUTDOWN")
	if err := server.Shutdown(exitCtx); err != nil {
		neo4jDriver.Close(exitCtx)
		mongoClient.Disconnect(exitCtx) //nolint:errcheck
		if db, err := gormDB.DB(); err != nil {
			log.Fatal(err)
		} else {
			db.Close()
		}
		log.Fatal("Server Shutdown:", err)
	}

	return nil
}
