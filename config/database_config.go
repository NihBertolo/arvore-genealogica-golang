package config

import "github.com/neo4j/neo4j-go-driver/neo4j"

func InitializeDatabase() neo4j.Driver {
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "golangNeo4j"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = false
	})

	if err != nil {
		panic("Falha ao conectar ao banco de dados: " + err.Error())
	}

	return driver
}
