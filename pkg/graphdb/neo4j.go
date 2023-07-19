package graphdb

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
)

func NewNeo4j(ctx context.Context, cfg configs.Neo4j) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(cfg.Uri, neo4j.BasicAuth(cfg.User, cfg.Password, ""))
	if err != nil {
		return nil, err
	}

	return driver, nil
}
