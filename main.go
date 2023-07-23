package main

import (
	_ "github.com/raaaaaaaay86/go-project-structure/docs" //nolint:typecheck
)

func main() {
	//err := Run(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}

//func Run(ctx context.Context) error {
//	// Config
//	config, err := configs.ReadYaml("./config/app_config.yaml")
//	if err != nil {
//		return err
//	}
//
//	// Logger
//	zapLogger, err := logger.NewZapLogger()
//	if err != nil {
//		return err
//	}
//	defer zapLogger.Sync() //nolint:errcheck
//
//	// Tracer
//	exporter, err := tracing.NewJaegerExporter(config.Jaeger.Endpoint)
//	if err != nil {
//		zapLogger.Error("jaeger exporter error", err)
//	}
//
//	appTracer, err := tracing.NewJaegerTracerProvider("application", exporter)
//	if err != nil {
//		zapLogger.Error("init application tracer failed.", err)
//	}
//
//	httpTracer, err := tracing.NewJaegerTracerProvider("http", exporter)
//	if err != nil {
//		zapLogger.Error("init http tracer failed", err)
//	}
//
//	repositoryTracer, err := tracing.NewJaegerTracerProvider("repository", exporter)
//	if err != nil {
//		zapLogger.Error("init repository tracer failed", err)
//	}
//
//	// Neo4j Driver
//	neo4jDriver, err := graphdb.NewNeo4j(ctx, config.Neo4j)
//	if err != nil {
//		return err
//	}
//	defer neo4jDriver.Close(ctx)
//
//	// Postgres Repositories
//	gormDB, err := gorm.NewPostgresConnection(config.Postgres)
//	if err != nil {
//		neo4jDriver.Close(ctx)
//		return err
//	}
//	userRepository := mongodb.NewUserRepository(repositoryTracer, gormDB)
//	videoPostRepository := mongodb.NewVideoPostRepository(repositoryTracer, gormDB, neo4jDriver)
//
//	// MongoDB Repositories
//	mongoClient, err := mongo.NewMongoDbConnection(config.MongoDB)
//	if err != nil {
//		neo4jDriver.Close(ctx)
//		mongoClient.Disconnect(ctx) //nolint:errcheck
//		return err
//	}
//	videoCommentRepository := mongodb.NewVideoCommentRepository(repositoryTracer, mongoClient)
//
//	// Helper Package
//	fileUploader := bucket.NewLocalUploader(config.BucketPath)
//	ffmpeg := convert.NewFfmpeg(config.BucketPath)
//
//	// Use Case
//	registerUseCase := auth.NewRegisterUserUseCase(appTracer, userRepository)
//	loginUseCase := auth.NewLoginUseCase(appTracer, userRepository)
//	uploadVideoUseCase := video.NewUploadVideoUseCase(appTracer, fileUploader, ffmpeg)
//	createVideoUseCase := video.NewCreateVideoUseCase(appTracer, videoPostRepository, userRepository)
//	createCommentUseCase := comment.NewCreateCommentUseCase(appTracer, videoCommentRepository, userRepository, videoPostRepository)
//	findCommentByVideoUseCase := comment.NewFindByVideoUseCase(appTracer, videoCommentRepository)
//	deleteCommentUseCase := comment.NewDeleteCommentUseCase(appTracer, videoCommentRepository)
//	forceDeleteCommentUseCase := comment.NewForceDeleteCommentUseCase(appTracer, videoCommentRepository)
//	likeVideoUseCase := video.NewLikeVideoUseCase(appTracer, videoPostRepository)
//	unlikeVideoUseCase := video.NewUnLikeVideoUseCase(appTracer, videoPostRepository)
//
//	// HTTP Server
//	authController := route.NewAuthenticationController(httpTracer, registerUseCase, loginUseCase)
//	videoController := route.NewVideoController(httpTracer, uploadVideoUseCase, createVideoUseCase, likeVideoUseCase, unlikeVideoUseCase)
//	commentController := route.NewCommentController(httpTracer, createCommentUseCase, findCommentByVideoUseCase, deleteCommentUseCase, forceDeleteCommentUseCase)
//
//	httpPort := fmt.Sprintf(":%d", config.Http.Port)
//	err = http.
//		NewHttpServer(authController, videoController, commentController).
//		Run(httpPort)
//	if err != nil {
//		neo4jDriver.Close(ctx)
//		mongoClient.Disconnect(ctx) //nolint:errcheck
//
//		if db, err := gormDB.DB(); err != nil {
//			return err
//		} else {
//			db.Close()
//		}
//
//		return err
//	}
//
//	return nil
//}
