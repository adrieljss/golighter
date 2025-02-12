package v1_api

import (
	"os"
	"testing"

	v1_auth "github.com/adrieljss/golighter/api/v1/auth"
	v1_health "github.com/adrieljss/golighter/api/v1/health"
	v1_users "github.com/adrieljss/golighter/api/v1/users"
	"github.com/adrieljss/golighter/platform"
)

var app *platform.Application

func TestMain(m *testing.M) {
	a := platform.App(true)
	app = &a
	defer app.CloseApp()
	SetupApiRoutes(app.FiberApp, app)
	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {
	v1_health.TestHealth(t, app)
}

func TestAuthFlow(t *testing.T) {
	v1_auth.TestAuthFlow(t, app)
}

func TestUsers(t *testing.T) {
	v1_users.TestUserPermissions(t, app)
}