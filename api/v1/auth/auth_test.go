package v1_auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/adrieljss/golighter/internal/testutils"
	"github.com/adrieljss/golighter/platform"
	"github.com/gofiber/fiber/v3"
)

func TestAuth(t *testing.T) {
	app := platform.App(true)
	defer app.CloseApp()
	SetupAuthRoutes(app.FiberApp, &app)

	{
		encoder := json.Marshal()
		testutils.TestRequests(t, app.FiberApp, []testutils.TestRequest{
			{
				Name:           "register_user_bad_request",
				Method:         http.MethodPost,
				Path:           "/api/v1/auth/register",
				Body: encoder.Encode(),
				ExpectedStatus: fiber.StatusBadRequest,
			},
		})
	}
}
