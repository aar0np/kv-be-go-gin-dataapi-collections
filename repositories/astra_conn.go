package repositories

import (
	"github.com/datastax/astra-db-go/astra"
	"github.com/datastax/astra-db-go/astra/options"
)

type AstraConfig struct {
	Token    string
	Keyspace string
	Endpoint string
}

func NewAstraSession(cfg AstraConfig) (*astra.Db, error) {

	client := astra.NewClient(
		options.API().SetToken(cfg.Token),
		options.API().SetKeyspace(cfg.Keyspace),
	)

	db := client.Database(cfg.Endpoint)

	return db, nil
}
