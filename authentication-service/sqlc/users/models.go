// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package users_sql

import ()

type User struct {
	ID              int32  `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	EncodedPassword string `json:"encodedPassword"`
}