package service

import (
	"time"
)

type CreateUserParams struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
	DOB       time.Time
}
