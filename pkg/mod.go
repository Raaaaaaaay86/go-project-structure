package pkg

import (
	"github.com/google/wire"
	"github.com/raaaaaaaay86/go-project-structure/pkg/bucket"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	convert "github.com/raaaaaaaay86/go-project-structure/pkg/convert/video"
	"github.com/raaaaaaaay86/go-project-structure/pkg/gorm"
	"github.com/raaaaaaaay86/go-project-structure/pkg/graphdb"
	"github.com/raaaaaaaay86/go-project-structure/pkg/logger"
	"github.com/raaaaaaaay86/go-project-structure/pkg/mongo"
	"github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
)

//const pkg = "pkg"

var PkgSet = wire.NewSet(
	configs.YamlConfigProvider,
	tracing.NewJaegerExporter,
	tracing.NewApplicationTracer,
	tracing.NewRepositoryTracer,
	tracing.NewHttpTracer,
	tracing.NewGormTracer,
	logger.NewZapLogger,
	graphdb.NewNeo4j,
	gorm.NewPostgresConnection,
	mongo.NewMongoDbConnection,
	bucket.NewLocalUploader,
	wire.Bind(new(bucket.Uploader), new(*bucket.LocalUploader)),
	convert.NewFfmpeg,
	wire.Bind(new(convert.IFfmpeg), new(*convert.Ffmpeg)),
)
