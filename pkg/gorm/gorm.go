package gorm

import (
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(c configs.Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Taipei",
		c.Host,
		c.User,
		c.Password,
		c.Schema,
		c.Port,
	)
	postgresConfig := postgres.Config{
		DSN: dsn,
	}
	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.New(postgresConfig), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
