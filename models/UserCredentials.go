package models

import (
	astratypes "github.com/datastax/astra-db-go/datatypes"
)

type UserCredentials struct {
	Email         string
	Password      string
	Userid        astratypes.UUID
	AccountLocked bool
}
