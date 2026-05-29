package models

import (
	"time"
)

type User struct {
	Userid        string    `json:"userId"`
	Email         string    `json:"email"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	AccountStatus string    `json:"accountStatus"`
	CreatedDate   time.Time `json:"createdAt"`
	LastLoginDate time.Time `json:"lastLoginDate"`
}
