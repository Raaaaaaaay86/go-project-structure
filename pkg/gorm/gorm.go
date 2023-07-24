package gorm

import (
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	tracingx "github.com/raaaaaaaay86/go-project-structure/pkg/tracing"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

func NewPostgresConnection(c *configs.YamlConfig, gormTracer tracingx.GormTracer) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Taipei",
		c.Postgres.Host,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Schema,
		c.Postgres.Port,
	)
	postgresConfig := postgres.Config{
		DSN: dsn,
	}
	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.New(postgresConfig), gormConfig)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(gormTracer)

	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}

	return db, nil
}
