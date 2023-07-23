package graphdb

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
)

func NewNeo4j(cfg *configs.YamlConfig) (neo4j.DriverWithContext, error) {
	uri := cfg.Neo4j.Uri
	user := cfg.Neo4j.User
	password := cfg.Neo4j.Password
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, err
	}

	return driver, nil
}
