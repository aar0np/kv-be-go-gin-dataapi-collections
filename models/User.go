package models

import (
	"time"

	astratypes "github.com/datastax/astra-db-go/datatypes"
)

type User struct {
	Userid        astratypes.UUID `json:"userId"`
	Email         string          `json:"email"`
	FirstName     string          `json:"firstName"`
	LastName      string          `json:"lastName"`
	AccountStatus string          `json:"accountStatus"`
	CreatedDate   time.Time       `json:"createdAt"`
	LastLoginDate time.Time       `json:"lastLoginDate"`
}
