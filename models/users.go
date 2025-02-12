package models

import "time"

type Permission int64 // not uint64 since postgresql cannot be used with unsigned

const (
	PermissionNone Permission = 0

	// permission to read other users' informations
	PermissionUsersRead Permission = (1 << 0) // 1

	// permission to edit other users' informations
	PermissionUsersWrite Permission = (1 << 1) // 2

	// permission to delete other users
	PermissionUsersDelete Permission = (1 << 2) // 4

	PermissionAllUsers = PermissionUsersRead | PermissionUsersWrite | PermissionUsersDelete
)

type User struct {
	UID          string     `json:"uid"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Permissions  Permission `json:"permissions"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (u *User) HasPermission(p Permission) bool {
	return (u.Permissions & p) == p
}

func (u *User) GrantPermission(p Permission) {
	u.Permissions |= p
}

func (u *User) RevokePermission(p Permission) {
	u.Permissions &^= p // equivalent to u.Permissions &= ~p
}
