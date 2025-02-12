package v1_users

import (
	"testing"

	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/stretchr/testify/assert"
)

func TestUserPermissions(t *testing.T, app *platform.Application) {
	t.Helper()

	t.Run("Grant Permissions", func(t *testing.T) {
		user := &models.User{Permissions: models.PermissionNone}
		user.GrantPermission(models.PermissionUsersRead | models.PermissionUsersWrite)
		assert.True(t, user.HasPermission(models.PermissionUsersRead))
		assert.True(t, user.HasPermission(models.PermissionUsersWrite))
		assert.False(t, user.HasPermission(models.PermissionUsersDelete))
	})

	t.Run("Revoke Permission", func(t *testing.T) {
		user := &models.User{Permissions: models.PermissionUsersRead | models.PermissionUsersWrite}
		user.RevokePermission(models.PermissionUsersWrite)
		assert.True(t, user.HasPermission(models.PermissionUsersRead))
		assert.False(t, user.HasPermission(models.PermissionUsersWrite))
	})

	t.Run("Revoke Multiple Permissions", func(t *testing.T) {
		user := &models.User{Permissions: models.PermissionAllUsers}
		user.RevokePermission(models.PermissionUsersRead | models.PermissionUsersWrite)
		assert.False(t, user.HasPermission(models.PermissionUsersRead))
		assert.False(t, user.HasPermission(models.PermissionUsersWrite))
		assert.True(t, user.HasPermission(models.PermissionUsersDelete))
	})
}
