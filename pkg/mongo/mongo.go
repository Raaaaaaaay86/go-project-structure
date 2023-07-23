package mongo

import (
	"context"
	"fmt"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDbConnection(c *configs.YamlConfig) (*mongo.Client, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/?authMechanism=SCRAM-SHA-256&tls=false",
		c.MongoDB.User,
		c.MongoDB.Password,
		c.MongoDB.Host,
		c.MongoDB.Port,
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}
