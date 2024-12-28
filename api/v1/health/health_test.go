package v1_health

import (
	"net/http"
	"testing"

	"github.com/adrieljss/golighter/internal/testutils"
	"github.com/adrieljss/golighter/platform"
)

func TestHealth(t *testing.T) {
	app := platform.App(true)
	defer app.CloseApp()
	SetupHealthRoutes(app.FiberApp)
	testutils.TestRequests(t, app.FiberApp, []testutils.TestRequest{
		{
			Name:           "health_check",
			Method:         http.MethodGet,
			Path:           "/health",
			ExpectedStatus: 200,
			ExpectedBody:   "OK",
		},
		{
			Name:           "route_not_found",
			Method:         http.MethodGet,
			Path:           "/not_found",
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   "Not Found",
		},
	})
}
