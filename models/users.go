package models

import "time"

const (
	UserRoleUser = 0
	// 1-98 reserved for other roles
	UserRoleAdmin = 99
)

type User struct {
	UID          string    `json:"uid"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         int       `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
