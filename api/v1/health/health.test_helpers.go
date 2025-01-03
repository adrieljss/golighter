package v1_health

import (
	"testing"

	"github.com/adrieljss/golighter/internal/testutils"
	"github.com/adrieljss/golighter/platform"
	"github.com/gofiber/fiber/v3"
)

func TestHealth(t *testing.T, app *platform.Application) {
	t.Run("Health Check", func(t *testing.T) {
		testutils.TestReqs(t, app.FiberApp, []testutils.TestRequest{
			{
				Method:         fiber.MethodGet,
				Path:           "/v1/health",
				ExpectedStatus: 200,
				ExpectedBody:   "OK",
			},
		})
	})
}
