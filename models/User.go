package models

import (
	"time"

	astradb "github.com/datastax/astra-db-go"
)

type User struct {
	Userid        astradb.timeuuid `json:"userId"`
	Email         string           `json:"email"`
	FirstName     string           `json:"firstName"`
	LastName      string           `json:"lastName"`
	AccountStatus string           `json:"accountStatus"`
	CreatedDate   time.Time        `json:"createdAt"`
	LastLoginDate time.Time        `json:"lastLoginDate"`
}

type DBUser struct {
	Userid         astradb.timeuuid `json:"userid"`
	Email          string           `json:"email"`
	FirstName      string           `json:"firstname"`
	LastName       string           `json:"lastname"`
	AccountStatus  string           `json:"account_status"`
	HashedPassword string 			`json:"hashed_password"`
	CreatedDate    time.Time        `json:"created_date"`
	LastLoginDate  time.Time        `json:"last_login_date"`
}