package models

import (
	"time"
)

type DBUser struct {
	Userid         string    `json:"userid"`
	Email          string    `json:"email"`
	FirstName      string    `json:"firstname"`
	LastName       string    `json:"lastname"`
	AccountStatus  string    `json:"account_status"`
	HashedPassword string    `json:"hashed_password"`
	CreatedDate    time.Time `json:"created_date"`
	LastLoginDate  time.Time `json:"last_login_date"`
}
