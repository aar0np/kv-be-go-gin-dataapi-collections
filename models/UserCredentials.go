package models

import (
	astradb "github.com/datastax/astra-db-go"
)

type UserCredentials struct {
	Email         string
	Password      string
	Userid        astradb.timeuuid
	AccountLocked bool
}