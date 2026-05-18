package repositories

import (
	"fmt"
	"path/filepath"
	astradb "github.com/datastax/astra-db-go"
    "github.com/datastax/astra-db-go/options"
)

type AstraConfig struct {
	Token    string
	Keyspace string
	Endpoint string
}

func NewAstraSession(cfg AstraConfig) (*apachegocql.Session, error) {

	client := astradb.NewClient(
		options.WithToken(cfg.Token)
		options.WithKeyspace(cfg.Keyspace)
	)

	db := client.Database(cfg.Endpoint)

	return db
}